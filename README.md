## aws-env

[![Build Status](https://travis-ci.com/telia-oss/aws-env.svg?branch=master)](https://travis-ci.com/telia-oss/aws-env)

A small library and binary for securely handling secrets in environment variables on AWS. Supports KMS, SSM Parameter store and secrets manager.

## Usage

Both the library and binary versions of `aws-env` will loop through the environment and exchange any variables prefixed with
`sm://`, `ssm://` and `kms://` with their secret value from Secrets manager, SSM Parameter store or KMS respectively. In order
for this to work, the instance profile (EC2), task role (ECS), or execution role (Lambda) must have the correct privileges in order
to retrive the secret values and/or decrypt the secret using KMS.

For instance:
- `export SECRETSMANAGER=sm://example/secrets-manager`
- `export PARAMETERSTORE=ssm:///example/parameter-store`
- `export KMSENCRYPTED=kms://<encrypted-secret>`

For information about which credentials are required for these actions:
- Secrets manager: `secretsmanager:GetSecretValue` on the resource. If the secret is encrypted with a non-default KMS key, it also requires `kms:Decrypt` on said key.
- SSM Parameter store: `ssm:GetParameter` on the resource. `kms:Decrypt` on the KMS key used to encrypt the secret.
- KMS: `kms:Decrypt` on the key used to encrypt the secret string.

#### Region

The region used is determined by looking for `AWS_DEFAULT_REGION` first, and if it is unset or empty, it will attempt to get the region from the EC2 Metadata endpoint.

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
	"github.com/aws/aws-sdk-go/aws/session"
	environment "github.com/telia-oss/aws-env"
)

func main() {
	// New AWS Session
	sess, err := session.NewSession()
	if err != nil {
		panic(fmt.Errorf("failed to create new aws session: %s", err))
	}

	// Populate secrets using aws-env
	env, err := environment.New(sess)
	if err != nil {
		panic(fmt.Errorf("failed to initialize aws-env: %s", err))
	}
	if err := env.Populate(); err != nil {
		panic(fmt.Errorf("failed to populate environment: %s", err))
	}
}
```
