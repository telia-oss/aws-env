package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/session"
	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	awsenv "github.com/telia-oss/aws-env"
)

const (
	cmdDelim = "--"
)

var command rootCommand

type rootCommand struct {
	Exec execCommand `command:"exec" description:"Execute a command."`
}

type execCommand struct {
	JSON bool `long:"json" description:"Use JSON formatting when logging."`
}

// Execute command
func (c *execCommand) Execute(args []string) error {
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

	logger := logrus.New()
	if c.JSON {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	env, err := awsenv.New(sess, logger)
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
