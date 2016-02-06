package main

import (
	"bytes"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/zeroturnaround/configo/exec"
	"github.com/zeroturnaround/configo/flatmap"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const envVariablePrefix = "CONFIGO_SOURCE_"
const configoPrefix = "CONFIGO:"

type env struct {
	key   string
	value string
}

type sourceConfig struct {
	priority int
	value    string
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

	environ := os.Environ()

	if err := resolveAll(environ); err != nil {
		panic(err)
	}

	if err := processTemplatedEnvs(environ); err != nil {
		panic(err)
	}

	exec.Execute(strings.Join(os.Args[1:], " "))
}

func resolveAll(environ []string) error {
	_, err := fromEnviron(environ).
		Where(func(kv T) (bool, error) { return strings.HasPrefix(kv.(env).key, envVariablePrefix), nil }).
		Select(func(kv T) (T, error) {
		priority, err := strconv.Atoi(strings.TrimLeft(kv.(env).key, envVariablePrefix))
		if err != nil {
			return nil, err
		}
		return sourceConfig{priority, kv.(env).value}, nil
	}).
		OrderBy(func(a T, b T) bool { return a.(sourceConfig).priority <= b.(sourceConfig).priority }).
		Select(func(pair T) (T, error) { return GetSource(pair.(sourceConfig).value) }).
		// Resolve in parallel because some sources might use IO and will take some time
		AsParallel().AsOrdered().
		Select(func(source T) (T, error) { return source.(Source).Get() }).
		Select(func(config T) (T, error) { return flatmap.Flatten(config.(map[string]interface{})), nil }).
		AsSequential().
		All(func(partialConfig T) (bool, error) {
		for key, value := range partialConfig.(map[string]interface{}) {
			os.Setenv(key, fmt.Sprintf("%v", value))
		}
		return true, nil
	})

	return err
}

func processTemplatedEnvs(environ []string) error {
	envMap := make(map[string]string)

	// Calculate fresh map of environment variables
	fromEnviron(os.Environ()).All(func(kv T) (bool, error) {
		envMap[kv.(env).key] = kv.(env).value
		return true, nil
	})

	_, err := fromEnviron(environ).
		Where(func(kv T) (bool, error) { return strings.HasPrefix(kv.(env).value, configoPrefix), nil }).
		All(func(kv T) (bool, error) {
		tmpl, err := template.New("").Parse(strings.TrimPrefix(kv.(env).value, configoPrefix))

		if err != nil {
			return false, err
		}

		var buffer bytes.Buffer
		if err = tmpl.Execute(&buffer, envMap); err != nil {
			return false, err
		}

		os.Setenv(kv.(env).key, buffer.String())
		return true, nil
	})

	return err
}

func fromEnviron(environ []string) Query {
	return From(environ).Select(func(kv T) (T, error) {
		pair := strings.SplitN(kv.(string), "=", 2)
		return env{pair[0], pair[1]}, nil
	})
}
