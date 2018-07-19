## aws-env

[![Build Status](https://travis-ci.com/telia-oss/aws-env.svg?branch=master)](https://travis-ci.com/telia-oss/aws-env)

A small library and binary for securely handling secrets in environment variables on AWS. Supports KMS, SSM Parameter store and secrets manager.

## Usage

Both the library and binary versions of `aws-env` will loop through the environment and exchange any variables prefixed with
`sm://`, `ssm://` and `kms://` with their secret value from Secrets manager, SSM Parameter store or KMS respectively. In order
for this to work, the instance profile (EC2), task role (ECS), or execution role (Lambda) must have the correct privileges in order
to retrive the secret values and/or decrypt the secret using KMS.

For instance:
- `SM=sm://example/secrets-manager`: Replaces the variable `SM` with the value of `example/secrets-manager` in Secrets manager.
- `SSM=ssm:///example/parameter-store`: Replaces the variable `SSM` with the value of `/example/parameter-store` in SSM Parameter store.
- `KMS=kms://<secret>`: Replaces the variable `KMS` with the KMS decrypted value of `<secret>`.

For information about which credentials are required for these actions:
- Secrets manager: `secretsmanager:GetSecretValue` on the resource. If the secret is encrypted with a non-default KMS key, it also requires `kms:Decrypt` on said key.
- SSM Parameter store: `ssm:GetParameter` on the resource. `kms:Decrypt` on the KMS key used to encrypt the secret.
- KMS: `kms:Decrypt` on the key used to encrypt the secret string.

#### Binary

Grab a binary from the [releases](https://github.com/telia-oss/aws-env/releases) and start your process with:

```bash
aws-env exec -- <command>
```

This will populate all the secrets in the environment, and hand over the process to your `<command>` with the same PID.

#### Library

Import the library and invoke it prior to parsing flags or reading environment variables:

```go
package main

import (
	"github.com/sirupsen/logrus"
	awsenv "github.com/telia-oss/aws-env"
)

func main() {
	// New logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	// New AWS Session
	sess, err := session.NewSession()
	if err != nil {
		logger.Fatalf("failed to create a new session: %s", err)
	}

	// Populate secrets using awsenv
	env, err := awsenv.New(sess, logger)
	if err != nil {
		logger.Fatalf("failed to initialize awsenv: %s", err)
	}
	if err := env.Replace(); err != nil {
		logger.Fatalf("failed to populate environment variables: %s", err)
	}
}
```
