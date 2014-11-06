package git

import (
	"fmt"
	"os/exec"
)

type ErrPathIsNotADirectory struct {
	path string
}

func (e ErrPathIsNotADirectory) Error() string {
	return fmt.Sprintf("The given path exists, and is not a directory: %s", e.path)
}

func (e ErrPathIsNotADirectory) Path() string {
	return e.Path()
}

type ErrDirectoryIsNotAGitRepository struct {
	path string
}

func (e ErrDirectoryIsNotAGitRepository) Error() string {
	return fmt.Sprintf("The given path exists, and is not a git repository: %s", e.path)
}

func (e ErrDirectoryIsNotAGitRepository) Path() string {
	return e.Path()
}

type ErrExecFailed struct {
	cmd    string
	err    *exec.ExitError
	stdErr []byte
}

func (e ErrExecFailed) Error() string {
	return fmt.Sprintf("Could not execute command '%s'. Standard error was:\n%s\n Underlying error was: %s", e.cmd, e.stdErr, e.err.Error())
}

func (e ErrExecFailed) Cmd() string {
	return e.cmd
}

func (e ErrExecFailed) StdErr() []byte {
	return e.stdErr
}

func (e ErrExecFailed) ExitError() *exec.ExitError {
	return e.err
}
