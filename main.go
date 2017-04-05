package main

import (
	"encoding/json"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/bsideup/configo/exec"
	"github.com/bsideup/configo/flatmap"
	"github.com/bsideup/configo/sources"
	"github.com/op/go-logging"
	"os"
	"strconv"
	"strings"
)

const envVariablePrefix = "CONFIGO_SOURCE_"

type env struct {
	key   string
	value string
}

type sourceWithPriority struct {
	priority int
	value    string
}

var log = logging.MustGetLogger("configo")
var loggingBackend = logging.NewLogBackend(os.Stdout, "", 0)

func main() {
	pattern := os.Getenv("CONFIGO_LOG_PATTERN")

	if len(pattern) == 0 {
		pattern = `%{color}%{time:15:04:05.999} [%{level:.1s}] %{message}%{color:reset}`
	}

	format := logging.MustStringFormatter(pattern)
	logging.SetBackend(logging.NewBackendFormatter(loggingBackend, format))

	levelString := os.Getenv("CONFIGO_LOG_LEVEL")

	var level logging.Level
	if len(levelString) > 0 {
		var err error
		level, err = logging.LogLevel(levelString)

		if err != nil {
			log.Warningf("%s", err)
		}
	} else {
		level = logging.WARNING
	}

	logging.SetLevel(level, "configo")

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%s", r)
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		panic("the required argument `command` was not provided")
	}

	environ := os.Environ()

	log.Debugf("Environment variables at start:\n\t%s", strings.Join(environ, "\n\t"))

	if err := resolveAll(environ); err != nil {
		panic(err)
	}

	if err := processTemplatedEnvs(environ); err != nil {
		panic(err)
	}

	exec.Execute(strings.Join(os.Args[1:], " "))
}

func resolveAll(environ []string) error {
	rawTSources, err := fromEnviron(environ).
		Where(func(kv T) (bool, error) { return strings.HasPrefix(kv.(env).key, envVariablePrefix), nil }).
		Select(func(kv T) (T, error) {
			priority, err := strconv.Atoi(strings.TrimLeft(kv.(env).key, envVariablePrefix))
			if err != nil {
				return nil, err
			}
			return sourceWithPriority{priority, kv.(env).value}, nil
		}).
		OrderBy(func(a T, b T) bool { return a.(sourceWithPriority).priority <= b.(sourceWithPriority).priority }).
		Select(func(it T) (T, error) {
			sourceBytes := []byte(it.(sourceWithPriority).value)
			rawSource := make(map[string]interface{})
			err := json.Unmarshal(sourceBytes, &rawSource)
			return rawSource, err
		}).
		Results()

	if err != nil {
		return err
	}

	rawSources := make([]map[string]interface{}, len(rawTSources))

	for k, v := range rawTSources {
		rawSources[k] = v.(map[string]interface{})
	}

	if len(rawSources) == 0 {
		log.Warning("No sources provided")
		return nil
	}

	if os.Getenv("CONFIGO_UPPERCASE_KEYS") == "0" {
		flatmap.UppercaseKeys = false
	}

	loader := sources.CompositeSource{
		Sources: rawSources,
	}

	resultEnv, err := loader.Get()

	if err != nil {
		return err
	}

	if len(resultEnv) == 0 {
		log.Info("No new env variables were added.")
		return nil
	}

	for key, value := range resultEnv {
		log.Infof("Setting %s to %v", key, value)
		os.Setenv(key, fmt.Sprintf("%v", value))
	}

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debugf("Environment variables after resolve:\n\t%s", strings.Join(os.Environ(), "\n\t"))
	}

	return nil
}

func fromEnviron(environ []string) Query {
	return From(environ).Select(func(kv T) (T, error) {
		pair := strings.SplitN(kv.(string), "=", 2)
		return env{pair[0], pair[1]}, nil
	})
}
