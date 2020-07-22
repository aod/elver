package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func EnvOrConfigContents(envar, file string) (io.Reader, error) {
	val, err := env(envar)
	if err == nil {
		return strings.NewReader(val), nil
	}

	baseConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(baseConfigDir, appName, file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func env(name string) (value string, err error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err = fmt.Errorf("no environment variable `%s` found", name)
	}

	return
}
