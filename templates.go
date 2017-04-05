package main

import (
	"bytes"
	"errors"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/bsideup/configo/parsers"
	"github.com/op/go-logging"
	"os"
	"strings"
	"text/template"
)

const configoPrefix = "CONFIGO:"

var customFuncs = func() template.FuncMap {
	result := template.FuncMap{
		"encrypt": encrypt,
		"decrypt": decrypt,
	}

	for _, el := range []string{"JSON", "YAML", "HCL", "TOML", "Properties"} {
		format := el
		result["from"+format] = func(source string) (map[string]interface{}, error) {
			parser := parsers.GetParser(format)
			if parser == nil {
				return nil, errors.New("Unknown format " + format)
			}
			result := make(map[string]interface{})
			return result, parser.Parse([]byte(source), result)
		}
	}

	return result
}()

func processTemplatedEnvs(environ []string) error {
	envMap := make(map[string]string)

	// Calculate fresh map of environment variables
	fromEnviron(os.Environ()).All(func(kv T) (bool, error) {
		envMap[kv.(env).key] = kv.(env).value
		return true, nil
	})

	count, err := fromEnviron(environ).
		Where(func(kv T) (bool, error) { return strings.HasPrefix(kv.(env).value, configoPrefix), nil }).
		CountBy(func(kv T) (bool, error) {
			tmpl, err := template.New(kv.(env).key).Funcs(customFuncs).Parse(strings.TrimPrefix(kv.(env).value, configoPrefix))

			if err != nil {
				return false, err
			}

			var buffer bytes.Buffer
			if err = tmpl.Execute(&buffer, envMap); err != nil {
				return false, err
			}

			key := kv.(env).key
			value := buffer.String()

			log.Infof("Setting templated variable `%s` to `%#v`", key, value)

			err = os.Setenv(key, value)

			if err != nil {
				return false, err
			}

			return true, nil
		})

	if err != nil {
		return err
	}

	if count > 0 {
		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("Environment variables after templates:\n\t%s", strings.Join(os.Environ(), "\n\t"))
		}
	}

	return nil
}
