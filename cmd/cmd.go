package cmd

import (
	"bytes"
	"fmt"

	"github.com/lefeck/ubuntu-autoinstaller/logger"

	"github.com/sirupsen/logrus"

	"os/exec"
	"strings"
	"time"
)

// CmdOptions holds optional command execution parameters.
type CmdOptions struct {
	Name string
}

// Options applies optional settings to CmdOptions.
type Options interface {
	Apply(*CmdOptions)
}

// ApplyOptions applies a list of Options to CmdOptions.
func (o *CmdOptions) ApplyOptions(opts []Options) {
	for _, opt := range opts {
		opt.Apply(o)
	}
}

// CmdName sets the command name string.
type CmdName string

// Apply sets the Name on CmdOptions.
func (c CmdName) Apply(opt *CmdOptions) {
	opt.Name = string(c)
}

// CmdExecutor defines a generic command executor interface.
type CmdExecutor interface {
	RunCmd(cmd interface{}, opts ...Options) (string, string, error)
	RunCmdWithAttempts(cmd interface{}, attempts int, timeout time.Duration, opts ...Options) (string, string, error)
}

// Executor is a default implementation of CmdExecutor.
type Executor struct {
}

// RunCmdWithAttempts runs command with attempts and timeout.
func (e *Executor) RunCmdWithAttempts(cmd interface{}, attempts int, timeout time.Duration, opts ...Options) (string, string, error) {
	options := &CmdOptions{}
	options.ApplyOptions(opts)

	ll := logger.Logger
	var (
		stdout string
		stderr string
		err    error
	)
	for i := 0; i < attempts; i++ {
		if stdout, stderr, err := e.RunCmd(cmd); err == nil {
			return stdout, stderr, err
		}
		ll.Warnf("Unable to execute cmd: %v. Attempt %d out of %d.", err, i, attempts)
		<-time.After(timeout)
	}
	errMsg := fmt.Errorf("failed to execute command after %d attempt, error: %v", attempts, err)
	return stdout, stderr, errMsg
}

// RunCmd runs a command once.
func (e *Executor) RunCmd(cmd interface{}, opts ...Options) (string, string, error) {
	options := &CmdOptions{}
	options.ApplyOptions(opts)

	if cmdstr, ok := cmd.(string); ok {
		return e.runCmdFromStr(cmdstr)
	}
	if cmdstr, ok := cmd.(*exec.Cmd); ok {
		return e.runCmdFromCmdObj(cmdstr)
	}
	return "", "", fmt.Errorf("could not interpret command from %v", cmd)
}

// runCmdFromStr runs command from a string.
func (e *Executor) runCmdFromStr(cmd string) (string, string, error) {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return "", "", fmt.Errorf("empty command string")
	}

	name := fields[0]
	args := fields[1:]
	execCmd := exec.Command(name, args...)

	return e.runCmdFromCmdObj(execCmd)
}

// runCmdFromCmdObj runs command from an exec.Cmd object.
func (e *Executor) runCmdFromCmdObj(cmd *exec.Cmd) (outStr string, errStr string, err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdStartTime := time.Now()
	err = cmd.Run()
	cmdDuration := time.Since(cmdStartTime)

	outStr, errStr = stdout.String(), stderr.String()

	// Only log command execution if enabled
	if logger.Config.ShowCommandOutput {
		// Construct log message based on command execution result, not stderr content
		// Some commands (like gpg) output normal information to stderr
		logFields := logrus.Fields{
			"cmd":         strings.Join(cmd.Args, " "),
			"duration":    cmdDuration.String(),
			"duration_ns": cmdDuration.Nanoseconds(),
		}

		if err != nil {
			// Command failed - log as ERROR
			logger.Logger.WithFields(logFields).
				Errorf("Command failed - stdout: %s, stderr: %s, error: %v", outStr, errStr, err)
		} else {
			// Command succeeded - log based on configured level
			// Include stderr if present (some commands output normal info to stderr)
			message := "Command executed successfully"
			if len(outStr) > 0 || len(errStr) > 0 {
				if len(errStr) > 0 {
					message = fmt.Sprintf("stdout: %s, stderr: %s", outStr, errStr)
				} else {
					message = fmt.Sprintf("stdout: %s", outStr)
				}
			}

			// Use configured log level for successful commands
			switch logger.Config.CommandLogLevel {
			case logrus.DebugLevel:
				logger.Logger.WithFields(logFields).Debug(message)
			case logrus.InfoLevel:
				logger.Logger.WithFields(logFields).Info(message)
			case logrus.WarnLevel:
				logger.Logger.WithFields(logFields).Warn(message)
			default:
				logger.Logger.WithFields(logFields).Info(message)
			}
		}
	}
	return outStr, errStr, err
}
