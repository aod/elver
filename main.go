package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"time"
)

var sessionID string

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sessionID = os.Getenv("AOC_SESSION")
	if len(sessionID) == 0 {
		fmt.Fprintln(os.Stderr, "no environment variable `AOC_SESSION` found")
		os.Exit(1)
	}

	err = runLatest(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runLatest(cwd string) error {
	years := aocYears()
	var year int
	var yPath string

	for i := len(years) - 1; i >= 0; i-- {
		p := path.Join(cwd, strconv.Itoa(years[i]))
		stat, err := os.Stat(p)

		if err != nil && os.IsNotExist(err) {
			continue
		} else if !stat.IsDir() {
			continue
		}

		year = years[i]
		yPath = p
		break
	}

	if len(yPath) == 0 {
		return errors.New("no advent year directory found")
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin")
	cmd.Dir = yPath
	var outb bytes.Buffer
	cmd.Stderr = &outb

	err := cmd.Run()
	if err != nil {
		errMsg := strings.Trim(outb.String(), "\n")
		return fmt.Errorf("%v: %s %s", err, yPath, errMsg)
	}

	p, err := plugin.Open(path.Join(yPath, fmt.Sprintf("%d.so", year)))
	if err != nil {
		return err
	}

	fmt.Println("AOC", year)

	for d := 25; d > 0; d-- {
		foundPart := false

	inner:
		for _, part := range [...]string{"A", "B"} {
			v, err := p.Lookup("Day" + strconv.Itoa(d) + part)
			if err != nil {
				break inner
			} else {
				foundPart = true
			}

			input, err := getInput(year, d)
			if err != nil {
				return err
			}

			solver := v.(func(string) (interface{}, error))
			start := time.Now()
			ans, err := solver(input)
			elapsed := time.Since(start)

			fmt.Printf("Day %d %s (%s):\n", d, part, elapsed)
			if err != nil {
				fmt.Println("[ERROR]", err)
			} else {
				fmt.Println(ans)
			}
		}

		if foundPart {
			return nil
		}
	}

	return errors.New("no solvers found")
}

func getInput(year int, day int) (string, error) {
	if err := validYear(year); err != nil {
		return "", err
	}
	if err := validDay(day); err != nil {
		return "", err
	}

	inputCacheDir, err := createCacheDir(year)
	if err != nil {
		return "", err
	}
	inputFile := filepath.Join(inputCacheDir, strconv.Itoa(day)+".txt")

	if _, err := os.Stat(inputFile); err != nil && os.IsNotExist(err) {
		req, err := createInputRequest(year, day, sessionID)
		if err != nil {
			return "", err
		}

		body, err := fetch(req)
		if err != nil {
			return "", err
		}

		if err := ioutil.WriteFile(inputFile, body, 0644); err != nil {
			return "", err
		}

		return string(body), nil
	}

	b, _ := ioutil.ReadFile(inputFile)
	return string(b), nil
}

func createCacheDir(year int) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheDir = filepath.Join(cacheDir, "aoc-inputs")
	inputFileDir := filepath.Join(cacheDir, strconv.Itoa(year))

	if err := os.MkdirAll(inputFileDir, 0744); err != nil {
		return "", err
	}
	return inputFileDir, nil
}

func createInputRequest(year, day int, sessionID string) (*http.Request, error) {
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

func validYear(year int) error {
	if year < 2015 {
		return errors.New("advent of code first started on 2015")
	}

	now := time.Now()
	if year > now.Year() {
		return fmt.Errorf("the latest advent of code is %d", now.Year()-1)
	}

	if year == now.Year() {
		nextAoc := time.Date(year, time.December, 1, 0, 0, 0, 0, time.FixedZone("EST/UTC-5", -5*60*60))

		if now.Before(nextAoc) {
			diff := nextAoc.Sub(now)
			return fmt.Errorf("advent of code %d starts in %s", year, diff)
		}
	}

	return nil
}

func aocYears() []int {
	curr := time.Now().Year()
	var years []int
	for y := 2015; y <= curr; y++ {
		years = append(years, y)
	}

	if validYear(years[len(years)-1]) != nil {
		years = years[:len(years)-1]
	}

	return years
}

func validDay(day int) error {
	if day < 1 || day > 25 {
		return fmt.Errorf("invalid day: %d", day)
	}

	return nil
}
