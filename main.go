package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/op/go-logging"
	"github.com/zeroturnaround/configo/exec"
	"github.com/zeroturnaround/configo/flatmap"
	"io"
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

type sourceContext struct {
	priority      int
	value         string
	loader        Source
	partialConfig map[string]interface{}
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
	count, err := fromEnviron(environ).
		Where(func(kv T) (bool, error) { return strings.HasPrefix(kv.(env).key, envVariablePrefix), nil }).
		Select(func(kv T) (T, error) {
		priority, err := strconv.Atoi(strings.TrimLeft(kv.(env).key, envVariablePrefix))
		if err != nil {
			return nil, err
		}
		return sourceContext{priority, kv.(env).value, nil, nil}, nil
	}).
		OrderBy(func(a T, b T) bool { return a.(sourceContext).priority <= b.(sourceContext).priority }).
		Select(func(context T) (T, error) {
		loader, err := GetSource(context.(sourceContext).value)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse source #%d: %s", context.(sourceContext).priority, err)
		}
		return sourceContext{context.(sourceContext).priority, context.(sourceContext).value, loader, nil}, nil
	}).
		// Resolve in parallel because some sources might use IO and will take some time
		AsParallel().AsOrdered().
		Select(func(context T) (T, error) {
		result, err := context.(sourceContext).loader.Get()

		if err != nil {
			return nil, fmt.Errorf("Failed to resolve source #%d: %s", context.(sourceContext).priority, err)
		}

		return sourceContext{context.(sourceContext).priority, context.(sourceContext).value, context.(sourceContext).loader, result}, nil
	}).
		AsSequential().
		CountBy(func(context T) (bool, error) {
		for key, value := range flatmap.Flatten(context.(sourceContext).partialConfig) {
			log.Infof("Source #%d: Setting %s to %v", context.(sourceContext).priority, key, value)
			os.Setenv(key, fmt.Sprintf("%v", value))
		}
		return true, nil
	})

	if err != nil {
		return err
	}

	if count == 0 {
		log.Warning("No sources provided")
	} else {
		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("Environment variables after resolve:\n\t%s", strings.Join(os.Environ(), "\n\t"))
		}
	}

	return nil
}

func processTemplatedEnvs(environ []string) error {
	envMap := make(map[string]string)

	// Calculate fresh map of environment variables
	fromEnviron(os.Environ()).All(func(kv T) (bool, error) {
		envMap[kv.(env).key] = kv.(env).value
		return true, nil
	})

	encryptionKey := os.Getenv("CONFIGO_ENCRYPTION_KEY")

	customFuncs := template.FuncMap{
		"encrypt": func(rawValue string) (string, error) {
			if len(encryptionKey) < 1 {
				return "", errors.New("CONFIGO_ENCRYPTION_KEY should be set in order to use `encrypt` function")
			}

			rawBytes := []byte(rawValue)

			if len(rawBytes)%aes.BlockSize != 0 {
				padding := aes.BlockSize - len(rawBytes)%aes.BlockSize
				padtext := bytes.Repeat([]byte{byte(0)}, padding)
				rawBytes = append(rawBytes, padtext...)
			}

			block, err := aes.NewCipher([]byte(encryptionKey))
			if err != nil {
				return "", err
			}
			ciphertext := make([]byte, aes.BlockSize+len(rawBytes))

			iv := ciphertext[:aes.BlockSize]
			if _, err := io.ReadFull(rand.Reader, iv); err != nil {
				return "", err
			}
			mode := cipher.NewCBCEncrypter(block, iv)
			mode.CryptBlocks(ciphertext[aes.BlockSize:], rawBytes)

			return base64.StdEncoding.EncodeToString(ciphertext), nil
		},
		"decrypt": func(encodedValue string) (string, error) {
			if len(encryptionKey) < 1 {
				return "", errors.New("CONFIGO_ENCRYPTION_KEY should be set in order to use `decrypt` function")
			}

			block, err := aes.NewCipher([]byte(encryptionKey))
			if err != nil {
				return "", err
			}

			b, err := base64.StdEncoding.DecodeString(encodedValue)
			if err != nil {
				return "", err
			}

			if len(b) < aes.BlockSize {
				return "", errors.New("ciphertext too short")
			}
			iv := b[:aes.BlockSize]
			b = b[aes.BlockSize:]

			if len(b)%aes.BlockSize != 0 {
				return "", errors.New("ciphertext is not a multiple of the block size")
			}

			mode := cipher.NewCBCDecrypter(block, iv)
			mode.CryptBlocks(b, b)

			b = bytes.TrimRight(b, "\x00")

			return string(b), nil
		},
	}

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

func fromEnviron(environ []string) Query {
	return From(environ).Select(func(kv T) (T, error) {
		pair := strings.SplitN(kv.(string), "=", 2)
		return env{pair[0], pair[1]}, nil
	})
}
