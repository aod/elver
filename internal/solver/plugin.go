package solver

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/aod/elver/aoc"
)

var (
	ErrSolverInvalidSignature = errors.New("invalid signature in plugin")
)

type Plugin = plugin.Plugin

func FromPlugin(p *Plugin, d aoc.Day, pt aoc.Part) (Func, error) {
	v, err := p.Lookup("Day" + d.String() + pt.String())
	if err != nil {
		return nil, fmt.Errorf("no solver found: %w", err)
	}
	solver, ok := v.(Func)
	if !ok {
		return nil, fmt.Errorf("incorrect func for Day%s%s got `%T`, expected `%T`: %w",
			d, pt, v, solver, ErrSolverInvalidSignature)
	}
	return solver, nil
}

func FromPluginBoth(p *Plugin, d aoc.Day) (Func, Func, error) {
	a, err := FromPlugin(p, d, aoc.Part1)
	if err != nil {
		return nil, nil, err
	}

	b, err := FromPlugin(p, d, aoc.Part2)
	if errors.Is(err, ErrSolverInvalidSignature) {
		return a, nil, err
	}

	return a, b, nil
}
