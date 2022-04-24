//go:build e2e

package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary() (string, error) {
	_, err := exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("go is not installed: %s", err)
	}
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", err)
	}
	out := filepath.Join(dir, "aws-env")

	build := exec.Command("go", "build", "-o", out, "-v", "main.go")
	if err := build.Run(); err != nil {
		return "", fmt.Errorf("failed to build binary: %s", err)
	}
	return out, nil
}

func TestMain(t *testing.T) {
	// Requires that the AWS CLI is installed
	_, err := exec.LookPath("aws")
	if err != nil {
		t.Fatalf("tests require that the aws cli is installed")
	}

	// Build the binary
	bin, err := buildBinary()
	if err != nil {
		t.Fatalf("failed to build binary for tests: %s", err)
	}
	defer os.RemoveAll(bin)

	tests := []struct {
		description string
		command     []string
		environment map[string]string
		expected    string
		shouldError bool
	}{
		{
			description: "basic binary works",
			environment: map[string]string{
				"AWS_DEFAULT_REGION": "eu-west-1",
				"SSM":                "ssm:///test/secret/ssm",
				"SM":                 "sm:///test/secret/sm",
				"KMS":                "kms://AQICAHhqWrYSJqgztV0Awtr37mVaA+xxzBYsa+JdcsrPJRazuAF3EZnHbZwp0w5KMmjkZ8EUAAAAZjBkBgkqhkiG9w0BBwagVzBVAgEAMFAGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMfDFmDeAvWVOaCA0kAgEQgCPn70De/Tfomg4iwq85K3QswDPZZ2xxWLTX7VCMJIp7+o2tHQ==",
			},
			command:     []string{"exec", "--", "sh", "-ce", "eval 'echo SSM=$SSM SM=$SM KMS=$KMS'"},
			expected:    "SSM=ssm-test SM=sm-test KMS=kms-test",
			shouldError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			cmd := exec.Command(bin, tc.command...)

			var env []string
			for k, v := range tc.environment {
				env = append(env, fmt.Sprintf("%s=%s", k, v))
			}
			for _, value := range os.Environ() {
				p := strings.SplitN(value, "=", 2)
				if _, exists := tc.environment[p[0]]; !exists {
					env = append(env, value)
				}
			}
			cmd.Env = env

			out, err := cmd.CombinedOutput()
			if err != nil && tc.shouldError == false {
				t.Fatalf("unexpected error: %s: output: %s", err, string(out))
			}
			if err == nil && tc.shouldError == true {
				t.Fatal("expected an error to occur")
			}
			if got := string(out); !strings.Contains(got, tc.expected) {
				t.Errorf("expected output to contain '%s' but got:\n%s\n", tc.expected, got)
			}
		})
	}
}
