package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func EnvOrConfigContents(envar, file string) (string, error) {
	val, err := env(envar)
	if err == nil {
		return val, err
	}

	baseConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(baseConfigDir, appName, file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func env(name string) (value string, err error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err = fmt.Errorf("no environment variable `%s` found", name)
	}

	return
}
