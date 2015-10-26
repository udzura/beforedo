package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func checkSuccessCommand(t *Task) (bool, error) {
	cmd := exec.Command("sh", "-c", t.SuccessCommand)
	out, err := cmd.CombinedOutput()

	var exitCode int
	var found bool
	if err == nil {
		exitCode = 0
		found = true
	} else {
		if exit, ok := err.(*exec.ExitError); ok {
			if s, ok := exit.Sys().(syscall.WaitStatus); ok {
				exitCode = s.ExitStatus()
				found = true
			}
		}
	}
	if !found && err != nil {
		return false, err
	}

	outStr := strings.TrimSpace(string(out))
	fmt.Printf("The check command %s exited: %d, output: `%s'\n", t.SuccessCommand, exitCode, outStr)

	return exitCode != 0, nil
}
