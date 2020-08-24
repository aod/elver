package aoc

import (
	"fmt"
	"net/http"
)

// CreateInputReq creates an HTTP request for retrieving the Advent of Code
// input given d.
func CreateInputReq(d Date, sessionID string) (*http.Request, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", d.Year, d.Day)

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
