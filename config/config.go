// Package config providies ways to retrieve configurations and such from the
// user's profile.
// To use this package correctly one must run `config.SetAppName` first which
// sets the global appName variable.
// This value is used to look up application-specific configurations, cache files, etc.
package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Contents returns an io.Reader with the given configuration file
// contents.
func Contents(file string) (io.Reader, error) {
	baseConfigDir, err := configDir()
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

func configDir() (string, error) {
	return os.UserConfigDir()
}
