package main

import (
	"flag"
	"fmt"
	"io"
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
		solversFinder = specificDaySolversFinder{day: day.Value}
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

	day, solvers, err := solversFinder.findSolvers(p)
	if err != nil {
		return err
	}

	input, err := getInput(year, day, sessionID)
	if err != nil {
		return err
	}

	stringInput := string(input)
	for i, solver := range solvers {
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
					if ans, err = solver(stringInput); err != nil {
						b.FailNow()
					}
				}
			})

			result.answer = ans
			result.err = err
			result.kind = resultKind{bench: &b}
		} else {
			start := time.Now()
			ans, err := solver(stringInput)
			elapsed := time.Since(start)

			result.answer = ans
			result.err = err
			result.kind = resultKind{normal: &elapsed}
		}

		printSolver(os.Stdout, result)
	}

	return nil
}

type solveResult struct {
	day    aoc.Day
	part   aoc.Part
	kind   resultKind
	err    error
	answer interface{}
}

type resultKind struct {
	bench  *testing.BenchmarkResult
	normal *time.Duration
}

func printSolver(w io.Writer, s solveResult) {
	fmt.Fprintf(w, "Day %d %v ", s.day, s.part)

	switch {
	case s.kind.bench != nil:
		b := s.kind.bench
		fmt.Fprintf(w, "(N=%d, %d ns/op, %d bytes/op, %d allocs/op):\n",
			b.N, b.NsPerOp(), b.AllocedBytesPerOp(), b.AllocsPerOp())
	case s.kind.normal != nil:
		fmt.Fprintf(w, "(%s):\n", s.kind.normal)
	}

	if s.err != nil {
		fmt.Fprintln(w, "[ERROR]", s.err)
	} else {
		fmt.Fprintln(w, s.answer)
	}
}
