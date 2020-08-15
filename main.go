package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/command"
	"github.com/aod/elver/config"
	"github.com/aod/elver/flags"
)

func main() {
	benchmarkFlag := flag.Bool("b", false, "enable benchmarking")

	year := &flags.IntRange{Value: 0, Min: int(aoc.FirstYear), Max: int(aoc.LastYear())}
	flag.Var(year, "y", "the `year` to run")

	day := &flags.IntRange{Value: 0, Min: int(aoc.FirstDay), Max: int(aoc.LastDay)}
	flag.Var(day, "d", "the `day` to run")

	flag.Parse()

	cwd, err := os.Getwd()
	handleError(err)

	config.SetAppName("elver")
	sessReader, err := config.EnvOrConfigContents("AOC_SESSION", "aoc_session")
	handleError(err)
	buf := new(strings.Builder)
	_, err = io.Copy(buf, sessReader)
	handleError(err)
	sessionID := strings.TrimSpace(buf.String())

	var dirFinder yearDirFinder = latestYearDirFinder{}
	if year.Value != 0 {
		dirFinder = specificYearDirFinder{year: aoc.Year(year.Value)}
	}

	var solversFinder solversFinder = latestSolversFinder{}
	if day.Value != 0 {
		solversFinder = specificDaySolversFinder{day: aoc.Day(day.Value)}
	}

	handleError(run(cwd, sessionID, *benchmarkFlag, dirFinder, solversFinder))
}

func run(cwd, sessionID string, benchmark bool, dirFinder yearDirFinder, solversFinder solversFinder) error {
	year, yPath, err := dirFinder.findYearDir(cwd)
	if err != nil {
		return err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	buildFile := filepath.Join(cacheDir, "elver", "builds", year.String())

	err = command.New("go build -buildmode=plugin -o=" + buildFile).Dir(yPath).Exec()
	if err != nil {
		return err
	}

	p, err := plugin.Open(buildFile)
	if err != nil {
		return err
	}

	day, solverA, solverB, err := solversFinder.findSolvers(p)
	if err != nil {
		return err
	}

	input, err := getInput(year, day, sessionID)
	if err != nil {
		return err
	}

	fmt.Println("AOC", year)

	stringInput := string(input)
	fmt.Fprintln(os.Stdout, solverA.solveResult(stringInput, benchmark))

	if solverB != nil {
		fmt.Fprintln(os.Stdout, solverB.solveResult(stringInput, benchmark))
	}

	return nil
}
