package config

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// EnvOrContents returns the environment variable envar or configuration
// contents of file.
func EnvOrContents(envar, file string) (io.Reader, error) {
	val, err := env(envar)
	if err == nil {
		return strings.NewReader(val), nil
	}
	return Contents(file)
}

func env(name string) (value string, err error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err = fmt.Errorf("no environment variable `%s` found", name)
	}
	return
}
