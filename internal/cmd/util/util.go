// Package util provides small utility functions.
package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Fetch executes req and returns the body and error if any.
// An error is returned when resp.StatusCode >= 400 or ioutil.Readall(resp.body)
// fails.
func Fetch(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// HandleError exits the program if err is not nil and prints it.
func HandleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
