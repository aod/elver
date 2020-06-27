package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"plugin"
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

func printSolver(day int, part string, solve func() (interface{}, error), benchmark bool) {
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
