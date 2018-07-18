package awsenv_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus/hooks/test"
	awsenv "github.com/telia-oss/aws-env"
	"github.com/telia-oss/aws-env/mocks"
)

func TestMain(t *testing.T) {
	tests := []struct {
		description string
		key         string
		value       string
		callsSM     bool
		smOutput    *secretsmanager.GetSecretValueOutput
		callsSSM    bool
		ssmOutput   *ssm.GetParameterOutput
		callsKMS    bool
		kmsOutput   *kms.DecryptOutput
	}{
		{
			description: "picks up sm secrets",
			key:         "TEST",
			value:       "sm://<secret-path>",
			callsSM:     true,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("secret"),
			},
		},
		{
			description: "picks up ssm secrets",
			key:         "TEST",
			value:       "ssm://<parameter-path>",
			callsSSM:    true,
			ssmOutput: &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String("secret"),
				},
			},
		},
		{
			description: "picks up kms secrets",
			key:         "TEST",
			value:       "kms://" + base64.StdEncoding.EncodeToString([]byte("<encrypted>")),
			callsKMS:    true,
			kmsOutput: &kms.DecryptOutput{
				Plaintext: []byte("secret"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := mocks.NewMockSMClient(ctrl)
			if tc.callsSM {
				sm.EXPECT().GetSecretValue(gomock.Any()).Times(1).Return(tc.smOutput, nil)
			}
			ssm := mocks.NewMockSSMClient(ctrl)
			if tc.callsSSM {
				ssm.EXPECT().GetParameter(gomock.Any()).Times(1).Return(tc.ssmOutput, nil)
			}
			kms := mocks.NewMockKMSClient(ctrl)
			if tc.callsKMS {
				kms.EXPECT().Decrypt(gomock.Any()).Times(1).Return(tc.kmsOutput, nil)
			}

			// Set environment
			old := os.Getenv(tc.key)
			if err := os.Setenv(tc.key, tc.value); err != nil {
				t.Fatalf("failed to set environment variable: %s", err)
			}
			// Set the old value before exiting
			defer func() {
				if err := os.Setenv(tc.key, old); err != nil {
					t.Fatalf("failed to set environment variable: %s", err)
				}
			}()

			// Run tests
			logger, _ := logrus.NewNullLogger()
			env := awsenv.NewTestManager(sm, ssm, kms, logger)
			if err := env.Replace(); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

		})
	}
}
