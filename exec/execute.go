package exec

import (
	"os"
	"syscall"
)

// Execute executes given command
func Execute(command string) {
	cmd := ShellInvocationCommand(command)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); cmd.Process == nil {
		panic(err)
	}
	os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}
