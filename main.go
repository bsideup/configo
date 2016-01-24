package main

import (
	"bytes"
	"fmt"
	"github.com/zeroturnaround/configo/flatmap"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"text/template"
)

const envVariablePrefix = "CONFIGO_SOURCE_"
const configoPrefix = "CONFIGO:"

type configResolveResult struct {
	config map[string]interface{}
	err    error
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Stderr.WriteString(fmt.Sprintln(r))
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		panic("the required argument `command` was not provided\n")
	}

	command := strings.Join(os.Args[1:], " ")

	sources := getSources()

	configs := getConfigs(sources)

	sourceKeys := make([]int, 0, len(configs))
	for key := range configs {
		sourceKeys = append(sourceKeys, key)
	}
	sort.Ints(sourceKeys)

	for _, sourceKey := range sourceKeys {
		result := configs[sourceKey]

		if result.err != nil {
			panic(result.err)
		}

		partialConfig := flatmap.Flatten(result.config)

		for key, value := range partialConfig {
			os.Setenv(key, fmt.Sprintf("%v", value))
		}
	}

	envVars := GetEnvironmentVariables()
	for key, value := range envVars {
		if strings.HasPrefix(value, configoPrefix) {
			tmpl, err := template.New(key).Parse(strings.TrimPrefix(value, configoPrefix))

			if err != nil {
				panic(err)
			}

			var buffer bytes.Buffer
			tmpl.Execute(&buffer, envVars)

			os.Setenv(key, buffer.String())
		}
	}

	execute(command)
}

func getSources() map[int]string {
	sources := make(map[int]string)
	for key, value := range GetEnvironmentVariables() {
		if strings.HasPrefix(key, envVariablePrefix) {
			index, err := strconv.Atoi(strings.TrimLeft(key, envVariablePrefix))

			if err != nil {
				panic(err)
			}

			sources[index] = value
		}
	}

	return sources
}

// Resolves the specified map of
func getConfigs(sources map[int]string) map[int]configResolveResult {
	configs := make(map[int]configResolveResult, len(sources))

	var waitGroup sync.WaitGroup
	for sourceKey := range sources {
		source := sources[sourceKey]

		waitGroup.Add(1)
		//TODO retry
		go func(sourceKey int, source string) {
			defer waitGroup.Done()

			config, err := GetConfig(source)

			configs[sourceKey] = configResolveResult{config, err}
		}(sourceKey, source)
	}
	waitGroup.Wait()

	return configs
}

func execute(command string) {
	cmd := ShellInvocationCommand(command)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); cmd.Process == nil {
		panic(err)
	}
	os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}
