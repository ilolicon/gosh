package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	DEFAULT_RUM_TIMEOUT = 3600
)

type Command struct {
	Cmd           string
	Timeout       time.Duration
	TerminateChan chan int
	Setpgid       bool
	command       *exec.Cmd
	stdout        bytes.Buffer
	stderr        bytes.Buffer
}

func NewCommand(cmd string, opts ...CommandOpt) *Command {
	c := &Command{
		Cmd:           cmd,
		Timeout:       DEFAULT_RUM_TIMEOUT * time.Second,
		TerminateChan: make(chan int),
		Setpgid:       false,
	}
	for _, opt := range opts {
		opt(c)
	}

	command := exec.Command("/bin/bash", "-c", c.Cmd)
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: c.Setpgid}
	command.Stderr = &c.stderr
	command.Stdout = &c.stdout
	c.command = command
	return c
}

func (c *Command) Run() error {
	if err := c.command.Start(); err != nil {
		return err
	}

	errorChan := make(chan error)
	go func() {
		defer close(errorChan)
		errorChan <- c.command.Wait()
	}()

	var err error
	select {
	case err = <-errorChan:
	case <-time.After(c.Timeout):
		if err = c.terminate(); err == nil {
			err = fmt.Errorf("cmd run timeout, cmd `%s`, time `%v`", c.Cmd, c.Timeout)
		}
	case <-c.TerminateChan:
		if err = c.terminate(); err == nil {
			err = fmt.Errorf("cmd is terminated, cmd `%s`", c.Cmd)
		}
	}
	return err
}

func (c *Command) Stderr() string {
	return strings.TrimSpace(c.stderr.String())
}

func (c *Command) Stdout() string {
	return strings.TrimSpace(c.stdout.String())
}

func (c *Command) terminate() error {
	if c.Setpgid {
		return syscall.Kill(-c.command.Process.Pid, syscall.SIGKILL)
	} else {
		return syscall.Kill(c.command.Process.Pid, syscall.SIGKILL)
	}
}

type CommandOpt func(*Command)

func WithTimeout(t time.Duration) CommandOpt {
	return func(c *Command) {
		c.Timeout = t
	}
}

func WithSetpgid(b bool) CommandOpt {
	return func(c *Command) {
		c.Setpgid = b
	}
}

func WithTerminateChan(ch chan int) CommandOpt {
	return func(c *Command) {
		c.TerminateChan = ch
	}
}
