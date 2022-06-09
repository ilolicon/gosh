package command

import (
	"testing"
)

func TestTaskRun(t *testing.T) {
	cmds := []string{
		"echo 'test task run'",
		"whoami",
		"date",
	}
	task := NewTask(cmds, 10)
	task.Run(false)
	if err := task.GetError(); err != nil {
		t.Errorf("cmd task running error: %s", err.Error())
	}
}
