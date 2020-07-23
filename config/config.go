package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ConfigContents(file string) (io.Reader, error) {
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
