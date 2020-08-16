package config

import (
	"os"
	"path/filepath"
)

var appName string

// SetAppName sets the global config app name which is used for retrieving
// specific configuration contents.
func SetAppName(name string) {
	appName = name
}

// Dir returns the user's config dir using appName as suffix.
func Dir() (string, error) {
	c, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(c, appName), nil
}

// CacheDir returns the user's cache dir using appName as suffix.
func CacheDir() (string, error) {
	c, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(c, appName), nil
}
