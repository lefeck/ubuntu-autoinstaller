package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test CmdOptions and Options interface integration.
func TestCmdOptionsApply(t *testing.T) {

	options := []Options{
		CmdName("echo hello"),
	}

	cmdOpts := &CmdOptions{}
	cmdOpts.ApplyOptions(options)

	assert.Equal(t, "echo hello", cmdOpts.Name)
}

// Test RunCmd with a string command.
func TestRunCmd_StringCommand(t *testing.T) {
	exec := &Executor{}
	stdout, stderr, err := exec.RunCmd("echo hello")

	assert.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
	assert.Equal(t, "", stderr)
}

// Test RunCmdWithAttempts succeeds on first try.
func TestRunCmdWithAttempts_Success(t *testing.T) {
	exec := &Executor{}
	stdout, stderr, err := exec.RunCmdWithAttempts("echo success", 3, time.Millisecond*10)

	assert.NoError(t, err)
	assert.Equal(t, "success\n", stdout)
	assert.Equal(t, "", stderr)
}

// Test RunCmdWithAttempts fails after retries.
func TestRunCmdWithAttempts_Failure(t *testing.T) {
	exec := &Executor{}
	// 故意执行一个不存在的命令
	stdout, stderr, err := exec.RunCmdWithAttempts("nonexistentcommand", 2, time.Millisecond*10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute command")
	assert.Equal(t, "", stdout)
	assert.Equal(t, "", stderr)
}

// Test RunCmd with an empty command.
func TestRunCmd_EmptyCommand(t *testing.T) {
	exec := &Executor{}
	stdout, stderr, err := exec.RunCmd("")

	assert.Error(t, err)
	assert.Equal(t, "empty command string", err.Error())
	assert.Equal(t, "", stdout)
	assert.Equal(t, "", stderr)
}
