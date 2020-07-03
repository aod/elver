package main

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/aod/elver/aoc"
)

var errInvalidSolverSignature = errors.New("invalid signature")

type solverFunc = func(string) (interface{}, error)
type solver struct {
	day   aoc.Day
	part  aoc.Part
	solve solverFunc
}

func pluginSolversAB(p *plugin.Plugin, day aoc.Day) (*solver, *solver, error) {
	sA, err := pluginSolverX(p, day, aoc.Part1)
	if err != nil {
		return nil, nil, fmt.Errorf("error in day %d: %w", day, err)
	}

	sB, err := pluginSolverX(p, day, aoc.Part2)
	if err != nil {
		return nil, nil, fmt.Errorf("error in day %d: %w", day, err)
	}

	return sA, sB, nil
}

func pluginSolverX(p *plugin.Plugin, day aoc.Day, part aoc.Part) (*solver, error) {
	v, err := p.Lookup(fmt.Sprintf("Day%d%s", day, part))
	if err != nil {
		return nil, fmt.Errorf("no solver found: %w", err)
	}

	solve, ok := v.(solverFunc)
	if !ok {
		return nil, fmt.Errorf("incorrect func Day%d%s got `%T`, expected `%T`: %w",
			day, part, v, solve, errInvalidSolverSignature)
	}

	return &solver{day, part, solve}, nil
}
