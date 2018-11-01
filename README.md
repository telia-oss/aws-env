## aws-env

[![Build Status](https://travis-ci.com/telia-oss/aws-env.svg?branch=master)](https://travis-ci.com/telia-oss/aws-env)

A small library and binary for securely handling secrets in environment variables on AWS. Supports KMS, SSM Parameter store and secrets manager. Inspired by [ssm-env](https://github.com/remind101/ssm-env).

## Usage

Both the library and binary versions of `aws-env` will loop through the environment and exchange any variables prefixed with
`sm://`, `ssm://` and `kms://` with their secret value from Secrets manager, SSM Parameter store or KMS respectively. In order
for this to work, the instance profile (EC2), task role (ECS), or execution role (Lambda) must have the correct privileges in order
to retrive the secret values and/or decrypt the secret using KMS.

For instance:
- `export SECRETSMANAGER=sm://<path>`
- `export PARAMETERSTORE=ssm://<path>`
- `export KMSENCRYPTED=kms://<encrypted-secret>`

Where `<path>` is the name of the secret in secrets manager or parameter store. `aws-env` will look up secrets in the region specified
in the `AWS_REGION` or `AWS_DEFAULT_REGION` environment variables, and if they are both unset/empty it will contact the EC2 Metadata endpoint 
(if possible) and use the region where it is deployed.

Required IAM privileges:
- Secrets manager: `secretsmanager:GetSecretValue` on the resource. And `kms:Decrypt` if not using the `aws/secretsmanager` key alias.
- SSM Parameter store: `ssm:GetParameter` on the resource. `kms:Decrypt` on the KMS key used to encrypt the secret.
- KMS: `kms:Decrypt` on the key used to encrypt the secret.

#### Binary

Grab a binary from the [releases](https://github.com/telia-oss/aws-env/releases) and start your process with:

```bash
aws-env exec -- <command>
```

This will populate all the secrets in the environment, and hand over the process to your `<command>` with the same PID. The
populated secrets are only made available to the `<command>` and 'disappear' when the process exits.

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

## Security

There are a couple of things to keep in mind when using `aws-env`:
- Spawned processes will inherit their parents environment by default. If your `<command>` spawns new processes they will inherit the environment _with the secrets already populated_, unless you hand-roll the environment for the new process.
- The environment for a running process can be read by the root user (and yourself) _after secrets have been populated_ by running `cat /proc/<pid>/environ` on Linux, and `ps eww <pid>` on OSX. However, if root or the spawning user is compromised a malicious user can just as easily fetch the secrets directly from the AWS API ¯\\_(ツ)_/¯
