package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/config"
	"github.com/aod/elver/internal/util"
)

func getInput(d aoc.Date, sessionID string) ([]byte, error) {
	inputCacheDir, err := createCacheDir(d.Year)
	if err != nil {
		return nil, err
	}
	inputFile := filepath.Join(inputCacheDir, d.Day.String()+".txt")

	if _, err := os.Stat(inputFile); err != nil && os.IsNotExist(err) {
		req, err := aoc.CreateInputReq(d, sessionID)
		if err != nil {
			return nil, err
		}

		body, err := util.Fetch(req)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(inputFile, body, 0644); err != nil {
			return nil, err
		}

		return body, nil
	}

	return ioutil.ReadFile(inputFile)
}

func createCacheDir(year aoc.Year) (string, error) {
	cacheDir, err := config.CacheDir()
	if err != nil {
		return "", err
	}

	inputFileDir := filepath.Join(cacheDir, "aoc-inputs", year.String())
	if err := os.MkdirAll(inputFileDir, 0744); err != nil {
		return "", err
	}
	return inputFileDir, nil
}
