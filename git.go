package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type gitCmd struct {
	path string
}

type Git interface {
	PathToGit(path string)
	SyncRepositoryToRemoteBranch(location string, repository string, ref string) error
}

func New() (Git, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	return &gitCmd{
		path,
	}, nil
}

func (g *gitCmd) PathToGit(path string) {
	g.path = path
}

func (g *gitCmd) SyncRepositoryToRemoteBranch(location string, repository string, ref string) error {
	s, err := os.Stat(location)

	if err == nil {
		if !s.IsDir() {
			return ErrPathIsNotADirectory{location}
		}

		// Rudimentary test to see if a directory is actually a git repository
		gitPath := path.Join(location, ".git")
		s, err = os.Stat(gitPath)
		if err == nil {
			if !s.IsDir() {
				return ErrPathIsNotADirectory{gitPath}
			}

			return g.syncRepository(location, repository, ref)
		}

		return ErrDirectoryIsNotAGitRepository{location}
	} else if os.IsNotExist(err) {
		parentDirectory := path.Dir(location)
		s, err = os.Stat(parentDirectory)

		if err == nil {
			if !s.IsDir() {
				return ErrPathIsNotADirectory{parentDirectory}
			}

			return g.cloneRepository(location, repository, ref)
		}

		return err
	}

	return err
}

func (g *gitCmd) cloneRepository(location string, repository string, ref string) error {
	_, err := g.execCommand("", "clone", repository, location)
	if err != nil {
		return err
	}

	_, err = g.execCommand(location, "fetch", "--tags")

	if err != nil {
		return err
	}

	_, err = g.execCommand(location, "checkout", ref)

	return err
}

func (g *gitCmd) execCommand(cwd string, args ...string) (string, error) {
	cmd := exec.Command(g.path, args...)
	cmd.Dir = cwd
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stderrReader)

	if err != nil {
		return "", err
	}
	stdErr := buf.Bytes()

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(stdoutReader)

	if err != nil {
		return "", err
	}

	stdOut := buf.String()

	if err = cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", ErrExecFailed{
				cmd:    strings.Join(cmd.Args, " "),
				err:    exitErr,
				stdErr: stdErr,
			}
		}
		return "", err
	}
	return strings.Trim(stdOut, " \n"), nil
}

func (g *gitCmd) syncRepository(location string, repository string, ref string) error {
	// We want to ensure that our origin is set up correctly
	// So we try to remove it, and for now assume that every ErrExecFailed error is because it did not exist
	_, err := g.execCommand(location, "remote", "rm", "origin")
	if err != nil {
		if _, ok := err.(ErrExecFailed); !ok {
			return err
		}
	}

	// Add the origin
	_, err = g.execCommand(location, "remote", "add", "origin", repository)
	if err != nil {
		return err
	}

	// Fetch it
	_, err = g.execCommand(location, "fetch", "origin")
	if err != nil {
		return err
	}

	// Check out desired ref
	_, err = g.execCommand(location, "checkout", ref)

	if err != nil {
		return err
	}

	// Now we need to merge with the remote, in case where our ref is a branch
	headBytes, err := ioutil.ReadFile(path.Join(location, ".git", "HEAD"))
	if err != nil {
		return err
	}

	head := fmt.Sprintf("%s", headBytes)

	isBranch := strings.HasPrefix(head, "ref: ")

	if isBranch {
		_, err = g.execCommand(location, "merge", "--ff-only", "origin/"+ref)
		return err
	}

	return nil
}
