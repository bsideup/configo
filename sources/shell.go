package sources

import (
	"github.com/bsideup/configo/exec"
	"github.com/bsideup/configo/parsers"
)

type ShellSource struct {
	Command string `json:"command"`
	Format  string `json:"format"`
}

func (shellSource *ShellSource) Get() (map[string]interface{}, error) {

	cmd := exec.ShellInvocationCommand(shellSource.Command)

	cmdOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	parsers.MustGetParser(shellSource.Format).Parse(cmdOut, result)

	return result, nil
}
