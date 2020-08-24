/*
Package command is a simple bare-bones wrapper around the std "os/exec" package.
The main feature is that the Exec function captures stderr output in the
returned error value.
*/
package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Interface uses a builder-like pattern to run programs via the terminal.
type Interface interface {
	Dir(string) Interface // Changes the current working directory.
	Exec() error          // Executes the command and returns an error if any.
}

type command struct {
	cmd exec.Cmd
}

func (c *command) Dir(dir string) Interface {
	c.cmd.Dir = dir
	return c
}

func (c *command) Exec() error {
	var errb bytes.Buffer
	c.cmd.Stderr = &errb

	err := c.cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %w",
			strings.Trim(errb.String(), "\n"), err)
	}

	return nil
}

// New returns a new command runner using cmd as the command to run.
func New(cmd string) Interface {
	cmdSplit := strings.Split(cmd, " ")
	name := cmdSplit[0]
	args := cmdSplit[1:]
	c := exec.Command(name, args...)

	return &command{*c}
}
