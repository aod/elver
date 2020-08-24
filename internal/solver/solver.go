package solver

import (
	"os"
	"testing"
	"time"

	"github.com/aod/elver/internal/util"

	"github.com/aod/elver/aoc"
)

type (
	Input  = string
	Output = interface{}
	Func   = func(Input) (Output, error)
)

type Solver struct {
	aoc.DatePart
	Solver Func
}

func (s Solver) Result(input string, rk ResultKind) Result {
	r := Result{DatePart: s.DatePart, Attr: ResultAttribute{ResultKind: rk}}
	switch rk {
	case BenchmarkResult:
		defer util.RedirectNull(&os.Stdout, &os.Stderr)()
		b := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if r.Answer, r.Err = s.Solve(input); r.Err != nil {
					b.FailNow()
				}
			}
		})
		r.Attr.B = &b
	case TimeResult:
		start := time.Now()
		r.Answer, r.Err = s.Solve(input)
		elapsed := time.Since(start)
		r.Attr.T = &elapsed
	}
	return r
}

func (s Solver) Solve(in Input) (Output, error) {
	return s.Solver(in)
}
