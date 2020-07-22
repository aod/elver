package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/command"
	"github.com/aod/elver/config"
	"github.com/aod/elver/flags"
)

func main() {
	benchmarkFlag := flag.Bool("b", false, "enable benchmarking")

	years := aoc.Years()
	year := &flags.IntRange{Value: 0, Min: aoc.FirstYear, Max: years[len(years)-1]}
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
		dirFinder = specificYearDirFinder{year: year.Value}
	}

	var solversFinder solversFinder = latestSolversFinder{}
	if day.Value != 0 {
		solversFinder = specificDaySolversFinder{day: aoc.Day(day.Value)}
	}

	handleError(runLatest(cwd, sessionID, *benchmarkFlag, dirFinder, solversFinder))
}

func runLatest(cwd, sessionID string, benchmark bool, dirFinder yearDirFinder, solversFinder solversFinder) error {
	year, yPath, err := dirFinder.findYearDir(cwd)
	if err != nil {
		return err
	}

	err = command.New("go build -buildmode=plugin").Dir(yPath).Exec()
	if err != nil {
		return err
	}

	p, err := plugin.Open(path.Join(yPath, fmt.Sprintf("%d.so", year)))
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

	stringInput := string(input)
	fmt.Fprintln(os.Stdout, solverA.solveResult(stringInput, benchmark))

	if solverB != nil {
		fmt.Fprintln(os.Stdout, solverB.solveResult(stringInput, benchmark))
	}

	return nil
}
