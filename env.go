package main

import (
	"github.com/subosito/gotenv"
	"os"
	"strings"
)

// GetEnvironmentVariables returns map of environment variables
func GetEnvironmentVariables() map[string]string {
	return gotenv.Parse(strings.NewReader(strings.Join(os.Environ(), "\n")))
}
