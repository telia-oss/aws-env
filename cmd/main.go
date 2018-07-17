package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/session"
	flags "github.com/jessevdk/go-flags"
	awsenv "github.com/telia-oss/aws-env"
)

const (
	cmdDelim = "--"
)

// RootCommand options
type rootCommand struct {
	Run runCommand `command:"run" description:"Run a command."`
}

// RunCommand options
type runCommand struct {
	Region string `long:"region" env:"REGION" description:"AWS region to use for API calls."`
}

// Execute command
func (c *runCommand) Execute(args []string) error {
	if len(args) <= 0 {
		return errors.New("please supply a command to run")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("failed to validate command: %s", err)
	}

	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create new session: %s", err)
	}

	env := awsenv.New(sess, c.Region)
	if err != nil {
		return fmt.Errorf("failed to create new manager: %s", err)
	}

	if err := env.Replace(); err != nil {
		return fmt.Errorf("failed to set up environment: %s", err)
	}

	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return fmt.Errorf("failed to execute command: %s", err)
	}
	return nil
}

// export SSM_TEST="ssm:///some/secret/ssm" SM_TEST="sm://some/secret/secretsmanager" SSM_TEST2="ssm:///some/secret2/ssm" SM_TEST2="sm:///some/secret2/secretsmanager" KMS_TEST="kms://AQICAHhCeTewA9s/tWLvjxRlvSrGuJ2Hx3m0oaBwrQrJOrGRKwEY1xstBFoqjvC9tFux0XndAAAAcjBwBgkqhkiG9w0BBwagYzBhAgEAMFwGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQM5HBtdBQV5709ed34AgEQgC+Fdrv/349TRP/WyeeY6Urxg7Y4+8aaa15BXrpJf514Ogn78V7ZxwjIWqZh786Llw=="
func main() {
	var command rootCommand
	_, err := flags.Parse(&command)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

}
