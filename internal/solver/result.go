package solver

import (
	"fmt"
	"testing"
	"time"

	"github.com/aod/elver/aoc"
)

type ResultKind int

const (
	BenchmarkResult ResultKind = iota
	TimeResult
)

type ResultAttribute struct {
	ResultKind
	B *testing.BenchmarkResult
	T *time.Duration
}

func (r ResultAttribute) String() string {
	switch r.ResultKind {
	case BenchmarkResult:
		return fmt.Sprintf("(N=%d, %d ns/op, %d bytes/op, %d allocs/op):\n",
			r.B.N, r.B.NsPerOp(), r.B.AllocedBytesPerOp(), r.B.AllocsPerOp())
	case TimeResult:
		return fmt.Sprintf("(%s):\n", r.T)
	}
	return ""
}

type Result struct {
	aoc.DatePart
	Attr   ResultAttribute
	Err    error
	Answer Output
}

func (s Result) String() string {
	res := fmt.Sprintf("Day %s%s ", s.Day, s.Part)
	res += s.Attr.String()
	res += s.answer()
	return res
}

func (s Result) answer() string {
	if err := s.Err; err != nil {
		return fmt.Sprintf("[ERROR] %s\n", err)
	}
	return fmt.Sprintf("%v", s.Answer)
}
