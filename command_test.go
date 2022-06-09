package command

import (
	"strings"
	"testing"
	"time"
)

func TestCmdRun(t *testing.T) {
	cmd := NewCommand("echo 'ilolicon.github.com'", WithTimeout(60*time.Second))
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	output := cmd.Stdout()
	if output != "ilolicon.github.com" {
		t.Errorf("cmd `echo ilolicon.github.com` output '%s' not eq 'syncd'", output)
	}
}

func TestCmdTimeout(t *testing.T) {
	cmd := NewCommand("sleep 2", WithTimeout(1*time.Second))
	err := cmd.Run()
	if err == nil {
		t.Error("cmd should run timeout and output error msg, but err is nil")
	}
	if !strings.Contains(err.Error(), "cmd run timeout") {
		t.Errorf("cmd run timeout output '%s' prefix not eq 'cmd run timeout'", err.Error())
	}
}

func TestCmdTerminate(t *testing.T) {
	terminateChan := make(chan int)
	var err error
	cmd := NewCommand("sleep 5", WithTerminateChan(terminateChan))
	go func() {
		err = cmd.Run()
		if err == nil {
			t.Error("cmd should be terminated and output error msg, but err is nil")
		}
		if !strings.Contains(err.Error(), "cmd is terminated") {
			t.Errorf("cmd is terminated output '%s' prefix not eq 'cmd is terminated'", err.Error())
		}
	}()

	terminateChan <- 1
}
