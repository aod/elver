package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Interface interface {
	Dir(string) Interface
	Exec() error
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

func New(cmd string) Interface {
	cmdSplit := strings.Split(cmd, " ")
	name := cmdSplit[0]
	args := cmdSplit[1:]
	c := exec.Command(name, args...)

	return &command{*c}
}
