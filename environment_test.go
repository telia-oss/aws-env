package environment_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	environment "github.com/telia-oss/aws-env"
	"github.com/telia-oss/aws-env/mocks"
)

func TestMain(t *testing.T) {
	tests := []struct {
		description string
		key         string
		value       string
		expect      string
		callsSM     bool
		smOutput    *secretsmanager.GetSecretValueOutput
		callsSSM    bool
		ssmOutput   *ssm.GetParameterOutput
		callsKMS    bool
		kmsOutput   *kms.DecryptOutput
	}{
		{
			description: "does not have sideffects for the regular environment",
			key:         "TEST",
			value:       "somevalue",
			expect:      "somevalue",
		},
		{
			description: "allows empty strings as secrets",
			key:         "TEST",
			value:       "sm://some/secret",
			expect:      "",
			callsSM:     true,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(""),
			},
		},
		{
			description: "picks up sm secrets",
			key:         "TEST",
			value:       "sm://some/secret",
			expect:      "secret",
			callsSM:     true,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("secret"),
			},
		},
		{
			description: "picks up ssm secrets",
			key:         "TEST",
			value:       "ssm://some/parameter",
			expect:      "secret",
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
			expect:      "secret",
			callsKMS:    true,
			kmsOutput: &kms.DecryptOutput{
				Plaintext: []byte("secret"),
			},
		},
		{
			description: "supports multi-value secrets in secrets manager",
			key:         "TEST",
			value:       "sm://some/secret#password",
			expect:      "secret",
			callsSM:     true,

			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(`{"name":"admin", "password":"secret"}`),
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
			env := environment.NewTestManager(sm, ssm, kms)
			if err := env.Populate(); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got, want := os.Getenv(tc.key), tc.expect; got != want {
				t.Errorf("\ngot: %s\nwanted: %s", got, want)
			}
		})
	}
}
