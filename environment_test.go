package environment_test

import (
	"encoding/base64"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"

	environment "github.com/telia-oss/aws-env"
	"github.com/telia-oss/aws-env/fakes"
)

func TestEnvironment(t *testing.T) {
	tests := []struct {
		description  string
		key          string
		value        string
		expect       string
		smCallCount  int
		smOutput     *secretsmanager.GetSecretValueOutput
		ssmCallCount int
		ssmOutput    *ssm.GetParameterOutput
		kmsCallCount int
		kmsOutput    *kms.DecryptOutput
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
			value:       "sm://<secret-path>",
			expect:      "",
			smCallCount: 1,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(""),
			},
		},
		{
			description: "picks up sm secrets",
			key:         "TEST",
			value:       "sm://<secret-path>",
			expect:      "secret",
			smCallCount: 1,
			smOutput: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("secret"),
			},
		},
		{
			description:  "picks up ssm secrets",
			key:          "TEST",
			value:        "ssm://<parameter-path>",
			expect:       "secret",
			ssmCallCount: 1,
			ssmOutput: &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String("secret"),
				},
			},
		},
		{
			description:  "picks up kms secrets",
			key:          "TEST",
			value:        "kms://" + base64.StdEncoding.EncodeToString([]byte("<encrypted>")),
			expect:       "secret",
			kmsCallCount: 1,
			kmsOutput: &kms.DecryptOutput{
				Plaintext: []byte("secret"),
			},
		},
		{
			description:  "allows JSON strings as secrets",
			key:          "TEST",
			value:        "ssm://<parameter-path>",
			expect:       `{"key":"secret"}`,
			ssmCallCount: 1,
			ssmOutput: &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String(`{"key":"secret"}`),
				},
			},
		},
		{
			description:  "supports multi-value secrets",
			key:          "TEST",
			value:        "ssm://<parameter-path>#key",
			expect:       "secret",
			ssmCallCount: 1,
			ssmOutput: &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String(`{"key":"secret"}`),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			fakeSM := &fakes.FakeSMClient{}
			fakeSM.GetSecretValueReturns(tc.smOutput, nil)

			fakeSSM := &fakes.FakeSSMClient{}
			fakeSSM.GetParameterReturns(tc.ssmOutput, nil)

			fakeKMS := &fakes.FakeKMSClient{}
			fakeKMS.DecryptReturns(tc.kmsOutput, nil)

			if err := os.Setenv(tc.key, tc.value); err != nil {
				t.Fatalf("failed to set environment variable: %s", err)
			}

			if err := environment.NewTestManager(fakeSM, fakeSSM, fakeKMS).Populate(); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			eq(t, tc.smCallCount, fakeSM.GetSecretValueCallCount())
			eq(t, tc.ssmCallCount, fakeSSM.GetParameterCallCount())
			eq(t, tc.kmsCallCount, fakeKMS.DecryptCallCount())
			eq(t, tc.expect, os.Getenv(tc.key))
		})
	}
}

func eq(t *testing.T, expected, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\nexpected:\n%v\n\ngot:\n%v", expected, got)
	}
}
