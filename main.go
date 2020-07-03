package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"plugin"
	"testing"
	"time"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/command"
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

	sessionID, err := env("AOC_SESSION")
	handleError(err)

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

	for i, solver := range [2]*solver{solverA, solverB} {
		if solver == nil {
			continue
		}

		// FIXME: This is dumb
		part := aoc.Part1
		if i == 1 {
			part = aoc.Part2
		}

		result := solveResult{
			day:  aoc.Day(day),
			part: part,
		}

		if benchmark {
			var ans interface{}
			var err error
			b := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if ans, err = solver.solve(stringInput); err != nil {
						b.FailNow()
					}
				}
			})

			result.answer = ans
			result.err = err
			result.kind = resultKind{bench: &b}
		} else {
			start := time.Now()
			ans, err := solver.solve(stringInput)
			elapsed := time.Since(start)

			result.answer = ans
			result.err = err
			result.kind = resultKind{normal: &elapsed}
		}

		fmt.Fprint(os.Stdout, result)
	}

	return nil
}
