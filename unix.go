// +build !windows

package main

import (
	"os/exec"
)

// ShellInvocationCommand creates exec.Cmd for UNIX-based platforms
func ShellInvocationCommand(args string) *exec.Cmd {
	return exec.Command("sh", "-c", args)
}
