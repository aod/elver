package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strconv"

	"github.com/aod/elver/aoc"
)

type yearDirFinder interface {
	findYearDir(string) (int, string, error)
}

type latestYearDirFinder struct{}

func (latestYearDirFinder) findYearDir(cwd string) (int, string, error) {
	years := aoc.Years()
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
	return 0, "", fmt.Errorf("no advent year %d directory found in %s", f.year, cwd)
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
