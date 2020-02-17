// common is the package of functionalities common to all commands
package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type ErrType uint8

const (
	// NotFoundErr represents the error for executing a command that is not
	// found within the path
	NotFoundErr ErrType = iota
	// RuntimeErr is an error triggered by the command
	RuntimeErr
	// UnexpectedErr is an error triggered by this program
	UnexpectedErr
)

// ExecErr is the error triggered by the execution of a command. It carry the
// ErrType and the error message
type ExecErr struct {
	Type ErrType
	msg  string
}

func (e *ExecErr) Error() string {
	return e.msg
}

var notFoundErr = &ExecErr{NotFoundErr, "command not found"}

const defaultFailedCode = -1

// ExecCmd executes a command and returns the exitcode as well as the stdout
// and stderr of the command executed. Returns the stdout or an ExecErr
// wrapping the stderr
func ExecCmd(name string, args ...string) (string, error) {
	// first we check if there is a git installation already
	stdout, stderr, err := runcmd(name, args...)
	if err != nil {
		log.WithError(err).WithField("name", name).Debugln("command triggered an error")
		exitCode := extractExitCode(err)
		if exitCode == defaultFailedCode {
			return "", notFoundErr
		}

		if len(stderr) > 0 {
			return "", &ExecErr{RuntimeErr, string(stderr)}
		}

		return "", &ExecErr{UnexpectedErr, err.Error()}
	}

	return string(stdout), nil
}

// extractExitCode from the error passed
func extractExitCode(err error) int {
	// base case
	if err == nil {
		return 0
	}

	// try to get the exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		return ws.ExitStatus()
	}

	// This will happen (in OSX) if `name` is not available in $PATH,
	// in this situation, exit code could not be get, and stderr will be
	// empty string very likely, so we use the default fail code, and format err
	// to string and set to stderr
	return defaultFailedCode
}

// runcmd executes a command and returns the stdout, stderr and an eventual
// error
func runcmd(name string, args ...string) ([]byte, []byte, error) {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString(" ")
	for i, arg := range args {
		if i != 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(arg)
	}

	log.Debugln(fmt.Sprintf("running cmd: `%s`", sb.String()))
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	if err := cmd.Run(); err != nil {
		return outbuf.Bytes(), errbuf.Bytes(), err
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}
