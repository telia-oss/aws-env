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

var command rootCommand

// RootCommand options
type rootCommand struct {
	Run    runCommand `command:"exec" description:"Run a command."`
	Region string     `long:"region" env:"AWS_DEFAULT_REGION" description:"AWS region to use for API calls."`
}

// RunCommand options
type runCommand struct{}

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

	env := awsenv.New(sess, command.Region)
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

func main() {
	_, err := flags.Parse(&command)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

}
