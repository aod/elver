/*
Package config providies ways to retrieve configurations and such from the
user's profile.

To use this package correctly one must run

	config.SetAppName("myappname")

first which sets the global appName variable.

This value is used to look up application-specific configurations, cache files, etc.
*/
package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
)

// Contents returns an io.Reader with the given configuration file
// contents.
func Contents(file string) (io.Reader, error) {
	baseConfigDir, err := Dir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(baseConfigDir, file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
