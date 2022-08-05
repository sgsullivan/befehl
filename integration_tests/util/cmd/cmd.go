package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Runs a system command, you must pass your cancel channel. Use this channel to tell it to
// kill the command. dir param represents the current working directory to run the command
// from. You can set to an empty string to be unset. evars refers to the environment variables
// to have set when the command runs.
func RunCmd(cmdName string, cmdArgs []string, cancel chan bool, dir string, evars []string) (stdout, stderr bytes.Buffer, err error) {
	cmd := exec.Command(cmdName, cmdArgs...) // #nosec
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if dir != "" {
		cmd.Dir = dir
	}

	env := os.Environ()
	env = append(env, evars...)

	cmd.Env = env

	cmdDone := make(chan error)
	if cancel == nil {
		cancel = make(chan bool)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if startErr := cmd.Start(); startErr != nil {
		err = fmt.Errorf("cmd.Start() returned an error: [%s]", startErr)
		return
	}

	go func() { cmdDone <- cmd.Wait() }()

WAIT:
	for {
		select {
		case <-cancel:
			killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			if killErr == nil {
				cmdDone <- fmt.Errorf("command [%s %s] was canceled; sent kill to cmd", cmdName, cmdArgs)
			} else {
				cmdDone <- fmt.Errorf("failed to .Kill() command [%s %s]: %s", cmdName, cmdArgs, killErr)
			}
		case cmdErr := <-cmdDone:
			err = cmdErr
			break WAIT
		}
	}

	return
}
