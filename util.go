package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func fetch(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func env(name string) (value string, err error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err = fmt.Errorf("no environment variable `%s` found", name)
	}
	return
}
