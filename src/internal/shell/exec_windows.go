//go:build windows

package shell

import (
	"os"
	"os/exec"
)

func replaceShell(shellPath string, env []string) error {
	cmd := exec.Command(shellPath)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	os.Exit(0)
	return nil // unreachable
}
