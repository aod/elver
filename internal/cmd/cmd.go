package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/aod/elver/internal/solver"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/command"
	"github.com/aod/elver/config"
	"github.com/aod/elver/flags"
	"github.com/aod/elver/internal/util"
)

// Execute is the entrypoint to elver.
func Execute(args []string) {
	benchmarkFlag := flag.Bool("b", false, "enable benchmarking")
	testFlag := flag.Bool("t", false, "enable testing")

	year := &flags.IntRange{Value: 0, Min: int(aoc.FirstYear), Max: int(aoc.LastYear())}
	flag.Var(year, "y", "the `year` to run")

	day := &flags.IntRange{Value: 0, Min: int(aoc.FirstDay), Max: int(aoc.LastDay)}
	flag.Var(day, "d", "the `day` to run")

	flag.Parse()

	cwd, err := os.Getwd()
	util.HandleError(err)

	config.SetAppName("elver")
	sessReader, err := config.EnvOrContents("AOC_SESSION", "aoc_session")
	util.HandleError(err)
	buf := new(strings.Builder)
	_, err = io.Copy(buf, sessReader)
	util.HandleError(err)
	sessionID := strings.TrimSpace(buf.String())

	var dirFinder yearDirFinder = latestYearDirFinder{}
	if year.Value != 0 {
		dirFinder = specificYearDirFinder{year: aoc.Year(year.Value)}
	}

	var solversFinder solversFinder = latestSolversFinder{}
	if day.Value != 0 {
		solversFinder = specificDaySolversFinder{day: aoc.Day(day.Value)}
	}

	opts := options{cwd, sessionID, *benchmarkFlag, *testFlag}
	util.HandleError(run(opts, dirFinder, solversFinder))
}

type options struct {
	cwd       string
	sessionID string
	benchmark bool
	test      bool
}

func run(opts options, dirFinder yearDirFinder, solversFinder solversFinder) error {
	year, yPath, err := dirFinder.findYearDir(opts.cwd)
	if err != nil {
		return err
	}

	cacheDir, err := config.CacheDir()
	if err != nil {
		return err
	}
	buildFile := filepath.Join(cacheDir, "builds", year.String())
	err = command.New("go build -buildmode=plugin -o=" + buildFile).Dir(yPath).Exec()
	if err != nil {
		return err
	}

	p, err := plugin.Open(buildFile)
	if err != nil {
		return err
	}

	day, funcA, funcB, err := solversFinder.findSolvers(p)
	if err != nil {
		return fmt.Errorf("%s: %w", year, err)
	}

	input, err := getInput(year, day, opts.sessionID)
	if err != nil {
		return err
	}

	k := solver.TimeResult
	if opts.benchmark {
		k = solver.BenchmarkResult
	}

	fmt.Println("AOC", year)
	stringInput := string(input)

	solverA := solver.Solver{
		DatePart: aoc.DatePart{
			Date: aoc.Date{Year: year, Day: day},
			Part: aoc.Part1,
		},
		Solver: funcA,
	}
	fmt.Fprintln(os.Stdout, solverA.Result(stringInput, k))

	if funcB == nil {
		return nil
	}

	solverB := solver.Solver{
		DatePart: aoc.DatePart{
			Date: aoc.Date{Year: year, Day: day},
			Part: aoc.Part2,
		},
		Solver: funcB,
	}
	fmt.Fprintln(os.Stdout, solverB.Result(stringInput, k))

	return nil
}
