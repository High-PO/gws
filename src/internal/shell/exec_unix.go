//go:build !windows

package shell

import "syscall"

func replaceShell(shellPath string, env []string) error {
	return syscall.Exec(shellPath, []string{shellPath}, env)
}
