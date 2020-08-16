package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/config"
)

func getInput(year aoc.Year, day aoc.Day, sessionID string) ([]byte, error) {
	inputCacheDir, err := createCacheDir(year)
	if err != nil {
		return nil, err
	}
	inputFile := filepath.Join(inputCacheDir, day.String()+".txt")

	if _, err := os.Stat(inputFile); err != nil && os.IsNotExist(err) {
		req, err := createInputRequest(year, day, sessionID)
		if err != nil {
			return nil, err
		}

		body, err := fetch(req)
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

func createInputRequest(year aoc.Year, day aoc.Day, sessionID string) (*http.Request, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{
		Name:   "session",
		Value:  sessionID,
		Domain: ".adventofcode.com",
		Path:   "/",
	})

	return req, nil
}
