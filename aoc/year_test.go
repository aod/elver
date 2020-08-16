package aoc

import (
	"reflect"
	"testing"
	"time"
)

func TestYears(t *testing.T) {
	testCases := []struct {
		t    time.Time
		want AdventYears
		desc string
	}{
		{
			t:    time.Date(2015, time.December, 1, 4, 59, 59, 0, time.UTC),
			want: AdventYears{},
			desc: "Pre AoC",
		},
		{
			t:    time.Date(2015, time.December, 1, 5, 0, 0, 1, time.UTC),
			want: AdventYears{2015},
			desc: "First AoC",
		},
		{
			t:    time.Date(2016, time.December, 1, 4, 59, 59, 0, time.UTC),
			want: AdventYears{2015},
			desc: "Close to 2nd",
		},
		{
			t:    time.Date(2016, time.December, 1, 5, 0, 0, 1, time.UTC),
			want: AdventYears{2015, 2016},
			desc: "2nd AoC",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			f := func() time.Time { return tc.t }
			got := years(f)
			ok := reflect.DeepEqual(got, tc.want)

			if !ok {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
