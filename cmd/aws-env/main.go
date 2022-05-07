package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jessevdk/go-flags"

	environment "github.com/telia-oss/aws-env"
)

var command rootCommand
var version string

type rootCommand struct {
	Version func()      `short:"v" long:"version" description:"Print the version and exit."`
	Exec    execCommand `command:"exec" description:"Execute a command."`
}

type execCommand struct{}

// Execute the exec subcommand.
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
		return fmt.Errorf("failed to create new aws session: %s", err)
	}

	env, err := environment.New(sess)
	if err != nil {
		return fmt.Errorf("failed to initialize aws-env: %s", err)
	}
	if err := env.Populate(); err != nil {
		return fmt.Errorf("failed to populate environment: %s", err)
	}

	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return fmt.Errorf("failed to execute command: %s", err)
	}
	return nil
}

func init() {
	command.Version = func() {
		fmt.Println(version)
		os.Exit(0)
	}
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
