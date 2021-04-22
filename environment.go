package environment

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	smPrefix     = "sm://"
	ssmPrefix    = "ssm://"
	kmsPrefix    = "kms://"
	envDelmiter  = "="
	mvsDelimiter = "#"
)

// Generate test fakes.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// SMClient (secrets manager client) for testing purposes.
//counterfeiter:generate -o ./fakes . SMClient
type SMClient interface {
	GetSecretValue(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

// SSMClient for testing purposes.
//counterfeiter:generate -o ./fakes . SSMClient
type SSMClient interface {
	GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
}

// KMSClient for testing purposes.
//counterfeiter:generate -o ./fakes . KMSClient
type KMSClient interface {
	Decrypt(*kms.DecryptInput) (*kms.DecryptOutput, error)
}

// NewTestManager for testing purposes.
func NewTestManager(sm SMClient, ssm SSMClient, kms KMSClient) *Manager {
	return &Manager{sm: sm, ssm: ssm, kms: kms}
}

// Manager handles API calls to AWS.
type Manager struct {
	sm  SMClient
	ssm SSMClient
	kms KMSClient
}

// New creates a new manager for populating secret values.
func New(sess *session.Session) (*Manager, error) {
	var config *aws.Config

	if os.Getenv("AWS_REGION") == "" && os.Getenv("AWS_DEFAULT_REGION") == "" {
		metadata := ec2metadata.New(sess)
		if !metadata.Available() {
			return nil, errors.New("'AWS_REGION' or 'AWS_DEFAULT_REGION' must be set when EC2 metadata is unavailable")
		}
		region, err := metadata.Region()
		if err != nil {
			return nil, fmt.Errorf("failed to get region from EC2 metadata: %s", err)
		}
		config = &aws.Config{Region: aws.String(region)}
	}

	return &Manager{
		sm:  secretsmanager.New(sess, config),
		ssm: ssm.New(sess, config),
		kms: kms.New(sess, config),
	}, nil
}

// Populate environment variables with their secret values from either Secrets manager, SSM Parameter store or KMS.
func (m *Manager) Populate() error {
	env := make(map[string]string)
	for _, v := range os.Environ() {
		var (
			found  bool
			secret string
			err    error
		)

		name, value, ok := strings.Cut(v, envDelmiter)
		if !ok {
			return fmt.Errorf("failed to parse environment variable with delimiter: %q", envDelmiter)
		}

		// # is not a legal character in secrets manager, parameter store or an
		// encrypted (and base64 encoded) string from KMS. I.e. it should only
		// be present if we are dealing with a multi-value secret.
		path, secretKey, isMultiValueSecret := strings.Cut(value, mvsDelimiter)

		if strings.HasPrefix(path, ssmPrefix) {
			secret, err = m.getParameter(strings.TrimPrefix(path, ssmPrefix))
			if err != nil {
				return fmt.Errorf("failed to get secret from parameter store: %q: %s", name, err)
			}
			found = true
		}
		if strings.HasPrefix(path, smPrefix) {
			secret, err = m.getSecretValue(strings.TrimPrefix(path, smPrefix))
			if err != nil {
				return fmt.Errorf("failed to get secret from secret manager: %q: %s", name, err)
			}
			found = true
		}
		if strings.HasPrefix(path, kmsPrefix) {
			secret, err = m.decrypt(strings.TrimPrefix(path, kmsPrefix))
			if err != nil {
				return fmt.Errorf("failed to decrypt kms secret: %q: %s", name, err)
			}
			found = true
		}
		if found {
			if isMultiValueSecret {
				o := make(map[string]string)
				if err := json.Unmarshal([]byte(secret), &o); err != nil {
					return fmt.Errorf("failed to unmarshal multi-value secret: %q", name)
				}

				secret, ok = o[secretKey]
				if !ok {
					return fmt.Errorf("failed to get multi-value secret with key (%q): %q", secretKey, name)
				}
			}
			env[name] = secret
		}
	}

	for name, secret := range env {
		if err := os.Setenv(name, secret); err != nil {
			return fmt.Errorf("failed to set environment variable: '%s': %s", name, err)
		}
	}
	return nil
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
