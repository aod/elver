package main

import (
	"errors"
	"fmt"
	"plugin"
	"testing"
	"time"

	"github.com/aod/elver/aoc"
)

var errInvalidSolverSignature = errors.New("invalid signature")

type solverFunc = func(string) (interface{}, error)
type solver struct {
	day   aoc.Day
	part  aoc.Part
	solve solverFunc
}

func (s solver) solveResult(input string, benchmark bool) solveResult {
	result := solveResult{
		day:  s.day,
		part: s.part,
	}

	if benchmark {
		var ans interface{}
		var err error

		b := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if ans, err = s.solve(input); err != nil {
					b.FailNow()
				}
			}
		})

		result.answer = ans
		result.err = err
		result.kind = resultKind{bench: &b}
	} else {
		start := time.Now()
		ans, err := s.solve(input)
		elapsed := time.Since(start)

		result.answer = ans
		result.err = err
		result.kind = resultKind{normal: &elapsed}
	}

	return result
}

func pluginSolversAB(p *plugin.Plugin, day aoc.Day) (*solver, *solver, error) {
	sA, err := pluginSolverX(p, day, aoc.Part1)
	if err != nil {
		return nil, nil, fmt.Errorf("error in day %d: %w", day, err)
	}

	sB, err := pluginSolverX(p, day, aoc.Part2)
	if errors.Is(err, errInvalidSolverSignature) {
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
