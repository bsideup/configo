// +build windows

package exec

import (
	"os/exec"
)

// ShellInvocationCommand creates exec.Cmd for Windows-based platforms
func ShellInvocationCommand(args string) *exec.Cmd {
	return exec.Command("cmd", "/C", args)
}
