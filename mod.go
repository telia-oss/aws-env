package awsenv

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/sirupsen/logrus"
)

const (
	smPrefix  = "sm://"
	ssmPrefix = "ssm://"
	kmsPrefix = "kms://"
)

// SMClient (secrets manager client) for testing purposes.
//go:generate mockgen -destination=mocks/mock_sm_client.go -package=mocks github.com/telia-oss/aws-env SMClient
type SMClient secretsmanageriface.SecretsManagerAPI

// SSMClient for testing purposes.
//go:generate mockgen -destination=mocks/mock_ssm_client.go -package=mocks github.com/telia-oss/aws-env SSMClient
type SSMClient ssmiface.SSMAPI

// KMSClient for testing purposes.
//go:generate mockgen -destination=mocks/mock_kms_client.go -package=mocks github.com/telia-oss/aws-env KMSClient
type KMSClient kmsiface.KMSAPI

// Manager handles API calls to AWS.
type Manager struct {
	sm     SMClient
	ssm    SSMClient
	kms    KMSClient
	logger *logrus.Logger
}

// New creates a new manager for handling AWS API calls.
func New(sess *session.Session, logger *logrus.Logger) (*Manager, error) {
	var (
		region string
		err    error
	)

	region = os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		metadata := ec2metadata.New(sess)
		if !metadata.Available() {
			return nil, errors.New("'AWS_DEFAULT_REGION' must be set when EC2 metadata is unavailable")
		}
		region, err = metadata.Region()
		if err != nil {
			return nil, fmt.Errorf("failed to get region from EC2 metadata: %s", err)
		}
	}

	config := &aws.Config{Region: aws.String("eu-west-1")}
	return &Manager{
		sm:     secretsmanager.New(sess, config),
		ssm:    ssm.New(sess, config),
		kms:    kms.New(sess, config),
		logger: logger,
	}, nil
}

// NewTestManager ...
func NewTestManager(sm SMClient, ssm SSMClient, kms KMSClient, logger *logrus.Logger) *Manager {
	return &Manager{sm: sm, ssm: ssm, kms: kms, logger: logger}
}

// Replace all environment variables with their secrets.
func (m *Manager) Replace() error {
	var errorCount int

	env := make(map[string]string)
	for _, v := range os.Environ() {
		var (
			secret string
			err    error
		)

		name, value := parseEnvironmentVariable(v)

		if strings.HasPrefix(value, ssmPrefix) {
			secret, err = m.getParameter(strings.TrimPrefix(value, ssmPrefix))
			if err != nil {
				err = fmt.Errorf("failed to get secret from parameter store: %s", err)
			}
		}
		if strings.HasPrefix(value, smPrefix) {
			secret, err = m.getSecretValue(strings.TrimPrefix(value, smPrefix))
			if err != nil {
				err = fmt.Errorf("failed to get secret from secret manager: %s", err)
			}
		}
		if strings.HasPrefix(value, kmsPrefix) {
			secret, err = m.decrypt(strings.TrimPrefix(value, kmsPrefix))
			if err != nil {
				err = fmt.Errorf("failed to decrypt kms secret: %s", err)
			}
		}

		if err != nil {
			if m.logger != nil {
				m.logger.WithField("variable", name).Warn(err)
			}
			errorCount++
		}
		if secret != "" {
			env[name] = secret
		}
	}

	for name, secret := range env {
		if err := os.Setenv(name, secret); err != nil {
			err = fmt.Errorf("failed to set environment variable: %s", err)
			if m.logger != nil {
				m.logger.WithField("variable", name).Warn(err)
			}
			errorCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("%d errors occured - check logs", errorCount)
	}
	return nil
}

func parseEnvironmentVariable(s string) (string, string) {
	pair := strings.SplitN(s, "=", 2)
	return pair[0], pair[1]
}

func (m *Manager) getSecretValue(path string) (out string, err error) {
	res, err := m.sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(path)})
	if err != nil {
		return "", err
	}

	if res.SecretString != nil {
		out = aws.StringValue(res.SecretString)
	} else {
		var data []byte
		if _, err := base64.StdEncoding.Decode(data, res.SecretBinary); err != nil {
			return "", fmt.Errorf("failed to decode binary secret: %s", err)
		}
		out = string(data)
	}
	return out, nil
}

func (m *Manager) getParameter(path string) (string, error) {
	res, err := m.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return aws.StringValue(res.Parameter.Value), nil
}

func (m *Manager) decrypt(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 cipher: %s", err)
	}
	res, err := m.kms.Decrypt(&kms.DecryptInput{CiphertextBlob: data})
	if err != nil {
		return "", err
	}
	return string(res.Plaintext), nil
}
