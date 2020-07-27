package aoc

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestEmptyDateRange(t *testing.T) {
	_, err := ParseDateRange("")
	want := ErrDateRangeEmpty
	if !errors.Is(err, want) {
		t.Errorf("got %q, want %s", err, want)
	}
}

func TestAllDateRange(t *testing.T) {
	tests := []struct {
		in      string
		nowTime time.Time
		outEnd  time.Time
	}{
		{
			in:      ":",
			nowTime: time.Date(2016, time.January, 1, 0, 0, 0, 0, Timezone),
			outEnd:  time.Date(2015, time.December, 25, 0, 0, 0, 0, Timezone),
		},
		{
			in:      ":",
			nowTime: time.Date(2016, time.December, 1, 0, 0, 0, 0, Timezone),
			outEnd:  time.Date(2016, time.December, 1, 0, 0, 0, 0, Timezone),
		},
		{
			in:      ":",
			nowTime: time.Date(2016, time.December, 21, 0, 0, 0, 0, Timezone),
			outEnd:  time.Date(2016, time.December, 21, 0, 0, 0, 0, Timezone),
		},
		{
			in:      ":",
			nowTime: time.Date(2016, time.December, 28, 0, 0, 0, 0, Timezone),
			outEnd:  time.Date(2016, time.December, 25, 0, 0, 0, 0, Timezone),
		},
	}

	for _, tt := range tests {
		t.Run(tt.outEnd.String(), func(t *testing.T) {
			nowf := func() time.Time { return tt.nowTime }
			v, err := parseDateRange(tt.in, nowf)

			if err != nil {
				t.Log(err)
			}

			if !tt.outEnd.Equal(v.End) {
				t.Errorf("got %q, want %q", v.End, tt.outEnd)
			}
		})
	}
}

func TestInvalidDateRange(t *testing.T) {
	tests := []struct {
		in  string
		err error
	}{
		{
			in:  "2014-1",
			err: ErrDateRangeMissingSeperator,
		},
		{
			in:  "2014-01:",
			err: ErrDateRangeDateNotAnAoC,
		},
		{
			in:  "2015-26:",
			err: ErrDateRangeDateNotAnAoC,
		},
		{
			in:  "2015-25:2019-31",
			err: ErrDateRangeDateNotAnAoC,
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			nowf := func() time.Time {
				return time.Date(2019, 1, 1, 0, 0, 0, 0, Timezone)
			}
			_, err := parseDateRange(tt.in, nowf)

			if !errors.Is(err, tt.err) {
				t.Errorf("got %q, want %q", err, tt.err)
			}
		})
	}
}


func TestDateRangeGen(t *testing.T) {
	tests := []struct {
		in   string
		want []time.Time
	}{
		{
			in: "2015-05:2015-10",
			want: []time.Time{
				time.Date(2015, 12, 5, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 6, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 7, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 8, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 9, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 10, 0, 0, 0, 0, Timezone),
			},
		},
		{
			in: "2015-10:2015-05",
			want: []time.Time{
				time.Date(2015, 12, 5, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 6, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 7, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 8, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 9, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 10, 0, 0, 0, 0, Timezone),
			},
		},
		{
			in: "2015-20:2016-05",
			want: []time.Time{
				time.Date(2015, 12, 20, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 21, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 22, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 23, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 24, 0, 0, 0, 0, Timezone),
				time.Date(2015, 12, 25, 0, 0, 0, 0, Timezone),
				time.Date(2016, 12, 1, 0, 0, 0, 0, Timezone),
				time.Date(2016, 12, 2, 0, 0, 0, 0, Timezone),
				time.Date(2016, 12, 3, 0, 0, 0, 0, Timezone),
				time.Date(2016, 12, 4, 0, 0, 0, 0, Timezone),
				time.Date(2016, 12, 5, 0, 0, 0, 0, Timezone),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			r, err := ParseDateRange(tt.in)
			if err != nil {
				t.Fatal(err)
			}

			c := r.Gen()
			if c == nil {
				t.Fatal("chan is nil")
			}

			var got []time.Time
			for v := range c {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v", got)
				t.Errorf("want: %v", tt.want)
				t.Errorf("%d != %d", len(got), len(tt.want))
			}
		})
	}
}
