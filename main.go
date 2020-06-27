package main

import (
	"bytes"
	"errors"
	"flag"
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
	"testing"
	"time"

	"github.com/aod/elver/flags"
)

func main() {
	benchmarkFlag := flag.Bool("b", false, "enable benchmarking")

	year := &flags.IntRange{Value: 0, Min: 2015, Max: 2019}
	flag.Var(year, "y", "the `year` to run")

	day := &flags.IntRange{Value: 0, Min: 1, Max: 25}
	flag.Var(day, "d", "the `day` to run")

	flag.Parse()

	cwd, err := os.Getwd()
	handleError(err)

	sessionID, err := env("AOC_SESSION")
	handleError(err)

	var dirFinder yearDirFinder = latestYearDirFinder{}
	if year.Value != 0 {
		dirFinder = specificYearDirFinder{year: year.Value}
	}

	var solversFinder solversFinder = latestSolversFinder{}
	if day.Value != 0 {
		solversFinder = specificDaySolversFinder{day: day.Value}
	}

	handleError(runLatest(cwd, sessionID, *benchmarkFlag, dirFinder, solversFinder))
}

type yearDirFinder interface {
	findYearDir(string) (int, string, error)
}

type latestYearDirFinder struct{}

func (latestYearDirFinder) findYearDir(cwd string) (int, string, error) {
	years := aocYears()
	for i := len(years) - 1; i >= 0; i-- {
		path := filepath.Join(cwd, strconv.Itoa(years[i]))

		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return years[i], path, nil
		}
	}

	return 0, "", fmt.Errorf("no advent year directory found in %s", cwd)
}

type specificYearDirFinder struct {
	year int
}

func (f specificYearDirFinder) findYearDir(cwd string) (int, string, error) {
	path := filepath.Join(cwd, strconv.Itoa(f.year))
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return f.year, path, nil
	}
	return 0, "", fmt.Errorf("advent year %d directory found in %s", f.year, cwd)
}

type solver = func(string) (interface{}, error)

type solversFinder interface {
	findSolvers(p *plugin.Plugin) (int, [2]solver, error)
}

type latestSolversFinder struct{}

func (latestSolversFinder) findSolvers(p *plugin.Plugin) (int, [2]solver, error) {
	var solvers [2]solver

	foundPart := false
	for day := 25; day > 0; day-- {
	inner:
		for i, part := range [...]string{"A", "B"} {
			v, err := p.Lookup("Day" + strconv.Itoa(day) + part)
			if err != nil {
				break inner
			}
			foundPart = true

			solver, ok := v.(func(string) (interface{}, error))
			if !ok {
				return 0, solvers, fmt.Errorf("found invalid solver signature for day %d: %T, expected: %T", day, v, solver)
			}

			solvers[i] = solver
		}

		if foundPart {
			return day, solvers, nil
		}
	}

	return 0, solvers, errors.New("no solvers found")
}

type specificDaySolversFinder struct {
	day int
}

func (f specificDaySolversFinder) findSolvers(p *plugin.Plugin) (int, [2]solver, error) {
	var solvers [2]solver

	for i, part := range [...]string{"A", "B"} {
		v, err := p.Lookup("Day" + strconv.Itoa(f.day) + part)
		if err != nil {
			if i != 0 {
				return f.day, solvers, nil
			}
			return 0, solvers, fmt.Errorf("no solvers found for day %d", f.day)
		}

		solver, ok := v.(func(string) (interface{}, error))
		if !ok {
			return 0, solvers, fmt.Errorf("found invalid solver signature for day %d: got %T, expected %T", f.day, v, solver)
		}

		solvers[i] = solver
	}

	return f.day, solvers, nil
}

func runLatest(cwd, sessionID string, benchmark bool, dirFinder yearDirFinder, solversFinder solversFinder) error {
	year, yPath, err := dirFinder.findYearDir(cwd)
	if err != nil {
		return err
	}

	err = buildPlugin(yPath)
	if err != nil {
		return err
	}

	p, err := plugin.Open(path.Join(yPath, fmt.Sprintf("%d.so", year)))
	if err != nil {
		return err
	}

	day, solvers, err := solversFinder.findSolvers(p)
	if err != nil {
		return err
	}

	input, err := getInput(year, day, sessionID)
	if err != nil {
		return err
	}
	stringInput := string(input)

	fmt.Println("AOC", year)
	for i, part := range [...]string{"A", "B"} {
		if solvers[i] != nil {
			solver := solvers[i]
			solve := func() (interface{}, error) {
				return solver(stringInput)
			}
			printSolver(day, part, solve, benchmark)
		}
	}

	return nil
}

type solverWrapper = func() (interface{}, error)

func printSolver(day int, part string, solve solverWrapper, benchmark bool) {
	fmt.Printf("Day %v %v ", day, part)

	var ans interface{}
	var err error

	if benchmark {
		b := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if ans, err = solve(); err != nil {
					b.FailNow()
				}
			}
		})

		fmt.Printf("(N=%d, %v ns/op, %v bytes/op, %v allocs/op):\n",
			b.N, b.NsPerOp(), b.AllocedBytesPerOp(), b.AllocsPerOp())
	} else {
		start := time.Now()
		ans, err = solve()
		elapsed := time.Since(start)

		fmt.Printf("(%s):\n", elapsed)
	}

	if err != nil {
		fmt.Println("[ERROR]", err)
	} else {
		fmt.Println(ans)
	}
}

func buildPlugin(dir string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin")
	cmd.Dir = dir
	var outb bytes.Buffer
	cmd.Stderr = &outb

	if err := cmd.Run(); err != nil {
		errMsg := strings.Trim(outb.String(), "\n")
		return fmt.Errorf("failed to run `%s` in %s, %s: %w", cmd, dir, errMsg, err)
	}

	return nil
}

func getInput(year, day int, sessionID string) ([]byte, error) {
	if err := validYear(year); err != nil {
		return nil, err
	}
	if err := validDay(day); err != nil {
		return nil, err
	}

	inputCacheDir, err := createCacheDir(year)
	if err != nil {
		return nil, err
	}
	inputFile := filepath.Join(inputCacheDir, strconv.Itoa(day)+".txt")

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
	years := make([]int, 0, curr-2015)
	for y := 2015; y < curr; y++ {
		years = append(years, y)
	}

	if validYear(curr) != nil {
		years = append(years, curr)
	}

	return years
}

func validDay(day int) error {
	if day < 1 || day > 25 {
		return fmt.Errorf("invalid day: %d", day)
	}

	return nil
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
