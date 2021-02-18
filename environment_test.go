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
		envKey      string
		secretKey   string
		envValue    string
		secretValue string
		expect      string
		json        bool
		callsSM     bool
		smOutput    *secretsmanager.GetSecretValueOutput
		callsSSM    bool
		ssmOutput   *ssm.GetParameterOutput
		callsKMS    bool
		kmsOutput   *kms.DecryptOutput
	}{
		{
			description: "does not have sideffects for the regular environment",
			envKey:      "TEST",
			envValue:    "somevalue",
			expect:      "somevalue",
		},
		{
			description: "allows empty strings as secrets",
			envKey:      "TEST",
			envValue:    "sm://<secret-path>",
			expect:      "",
			callsSM:     true,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(""),
			},
		},
		{
			description: "picks up sm secrets",
			envKey:      "TEST",
			envValue:    "sm://<secret-path>",
			expect:      "secret",
			callsSM:     true,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("secret"),
			},
		},
		{
			description: "picks up ssm secrets",
			envKey:      "TEST",
			envValue:    "ssm://<parameter-path>",
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
			envKey:      "TEST",
			envValue:    "kms://" + base64.StdEncoding.EncodeToString([]byte("<encrypted>")),
			expect:      "secret",
			callsKMS:    true,
			kmsOutput: &kms.DecryptOutput{
				Plaintext: []byte("secret"),
			},
		},
		{
			description: "test multi key-value secret within json",
			envKey:      "TEST",
			envValue:    "sm://<secret-path>",
			expect:      "sm://<secret-path>",
			secretKey:   "password",
			secretValue: "secret",
			json:        true,
			callsSM:     true,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("{\"password\": \"secret\"}"),
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
			old := os.Getenv(tc.envKey)
			if err := os.Setenv(tc.envKey, tc.envValue); err != nil {
				t.Fatalf("failed to set environment variable: %s", err)
			}
			// Set the old value before exiting
			defer func() {
				if err := os.Setenv(tc.envKey, old); err != nil {
					t.Fatalf("failed to set environment variable: %s", err)
				}
			}()

			// Run tests
			env := environment.NewTestManager(sm, ssm, kms)
			if err := env.Populate(); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if tc.json != true {
				if got, want := os.Getenv(tc.envKey), tc.expect; got != want {
					t.Errorf("\ngot: %s\nwanted: %s", got, want)
				}
			} else {
				if got, want := os.Getenv(tc.envKey), tc.expect; got != want {
					t.Errorf("\ngot: %s\nwanted: %s", got, want)
				}
				if got, want := os.Getenv(tc.secretKey), tc.secretValue; got != want {
					t.Errorf("\ngot: %s\nwanted: %s", got, want)
				}
			}
		})
	}
}
