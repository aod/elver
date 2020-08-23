package cmd

import (
	"errors"
	"fmt"
	"plugin"
	"sort"

	"github.com/aod/elver/aoc"
	"github.com/aod/elver/internal/solver"
)

type yearDirFinder interface {
	findYearDir(string) (aoc.Year, string, error)
}

type latestYearDirFinder struct{}

func (latestYearDirFinder) findYearDir(cwd string) (aoc.Year, string, error) {
	years := aoc.Years()
	sort.Sort(sort.Reverse(years))
	year, path, err := years.FirstYearDir(cwd)
	if err != nil {
		return 0, "", fmt.Errorf("no advent year directory found in %s: %w", cwd, err)
	}
	return year, path, nil
}

type specificYearDirFinder struct{ year aoc.Year }

func (f specificYearDirFinder) findYearDir(cwd string) (aoc.Year, string, error) {
	p, err := f.year.FindDir(cwd)
	if err != nil {
		return 0, "", fmt.Errorf("no advent year %d directory found in %s: %w", f.year, cwd, err)
	}
	return f.year, p, nil
}

type solversFinder interface {
	findSolvers(p *plugin.Plugin) (aoc.Day, solver.Func, solver.Func, error)
}

type latestSolversFinder struct{}

func (latestSolversFinder) findSolvers(p *plugin.Plugin) (aoc.Day, solver.Func, solver.Func, error) {
	for day := aoc.LastDay; day >= aoc.FirstDay; day-- {
		a, err := solver.FromPlugin(p, day, aoc.Part1)
		if errors.Is(err, solver.ErrSolverInvalidSignature) {
			return day, nil, nil, err
		} else if err != nil {
			continue
		}

		b, err := solver.FromPlugin(p, day, aoc.Part2)
		if errors.Is(err, solver.ErrSolverInvalidSignature) {
			return day, nil, nil, err
		}

		return day, a, b, nil
	}

	return 0, nil, nil, fmt.Errorf("no solvers found")
}

type specificDaySolversFinder struct{ day aoc.Day }

func (f specificDaySolversFinder) findSolvers(p *plugin.Plugin) (aoc.Day, solver.Func, solver.Func, error) {
	a, err := solver.FromPlugin(p, f.day, aoc.Part1)
	if err != nil {
		return f.day, nil, nil, err
	}

	b, err := solver.FromPlugin(p, f.day, aoc.Part2)
	if errors.Is(err, solver.ErrSolverInvalidSignature) {
		return f.day, nil, nil, err
	}

	return f.day, a, b, nil
}
