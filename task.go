package command

import (
	"fmt"
	"time"
)

type Task struct {
	Commands []string
	done     bool
	timeout  int
	termChan chan int
	err      error
	result   []*TaskResult
}

type TaskResult struct {
	Cmd     string `json:"cmd"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Success bool   `json:"success"`
}

func NewTask(cmds []string, timeout int) *Task {
	return &Task{
		Commands: cmds,
		termChan: make(chan int),
		timeout:  timeout,
	}
}

func (t *Task) Run(errExit bool) {
	for _, cmd := range t.Commands {
		result, err := t.next(cmd)
		t.result = append(t.result, result)
		if err != nil {
			t.err = fmt.Errorf("task run failed, %w", err)
			if errExit {
				break
			}
		}
	}
	t.done = true
}

func (t *Task) GetError() error {
	return t.err
}

func (t *Task) Result() []*TaskResult {
	return t.result
}

func (t *Task) Terminate() {
	if !t.done {
		t.termChan <- 1
	}
}

func (t *Task) next(cmd string) (*TaskResult, error) {
	result := &TaskResult{Cmd: cmd}
	command := NewCommand(cmd,
		WithTimeout(time.Duration(t.timeout)*time.Second),
		WithTerminateChan(t.termChan),
	)
	if err := command.Run(); err != nil {
		result.Stderr = command.Stderr()
		return result, err
	}
	result.Stdout = command.Stdout()
	result.Success = true
	return result, nil
}
