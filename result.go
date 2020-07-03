package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/aod/elver/aoc"
)

type resultKind struct {
	bench  *testing.BenchmarkResult
	normal *time.Duration
}

type solveResult struct {
	day    aoc.Day
	part   aoc.Part
	kind   resultKind
	err    error
	answer interface{}
}

func (s solveResult) String() string {
	result := fmt.Sprintf("Day %d %v ", s.day, s.part)

	switch {
	case s.kind.bench != nil:
		b := s.kind.bench
		result += fmt.Sprintf("(N=%d, %d ns/op, %d bytes/op, %d allocs/op):\n",
			b.N, b.NsPerOp(), b.AllocedBytesPerOp(), b.AllocsPerOp())
	case s.kind.normal != nil:
		result += fmt.Sprintf("(%s):\n", s.kind.normal)
	}

	if s.err != nil {
		result += fmt.Sprintf("[ERROR] %s\n", s.err)
	} else {
		result += fmt.Sprintf("%v", s.answer)
	}

	return result
}
