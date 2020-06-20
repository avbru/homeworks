package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(task []string, env Environment) (returnCode int) {
	cmd := exec.Command(task[0], task[1:]...) //nolint:gosec
	cmd.Stdout, cmd.Stdin, cmd.Stderr = os.Stdout, os.Stdin, os.Stderr

	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			fmt.Printf("cannot set env: %s, reason: %s\n", k, err)
		}
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("error executing command: %s, reason: %s\n", task[0], err)
	}

	return cmd.ProcessState.ExitCode()
}
