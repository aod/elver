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

type solversFinder interface {
	findSolvers(p *plugin.Plugin) (aoc.Day, *solver, *solver, error)
}

type latestSolversFinder struct{}

func (latestSolversFinder) findSolvers(p *plugin.Plugin) (aoc.Day, *solver, *solver, error) {
	for day := aoc.LastDay; day >= aoc.FirstDay; day-- {
		sA, sB, err := pluginSolversAB(p, day)

		if errors.Is(err, errInvalidSolverSignature) {
			return day, nil, nil, err
		}
		if err == nil {
			return day, sA, sB, err
		}
	}

	return 0, nil, nil, fmt.Errorf("no solvers found")
}

type specificDaySolversFinder struct {
	day aoc.Day
}

func (f specificDaySolversFinder) findSolvers(p *plugin.Plugin) (aoc.Day, *solver, *solver, error) {
	sA, sB, err := pluginSolversAB(p, f.day)
	if err != nil {
		return 0, nil, nil, err
	}

	return f.day, sA, sB, err
}
