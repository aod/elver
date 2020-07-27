package aoc

import (
	"errors"
	"strings"
	"time"
)

const dateLayout = "2006-02"

type DateRange struct {
	Start, End time.Time
	err        error
}

func (d DateRange) Gen() <-chan time.Time {
	out := make(chan time.Time)

	go func() {
		defer close(out)

		if d.err != nil {
			return
		}

		t := d.Start
		for t.Before(d.End) || t.Equal(d.End) {
			out <- t

			if t.Day() >= 25 {
				t = time.Date(t.Year()+1, time.December, int(FirstDay), 0, 0, 0, 0, Timezone)
			} else {
				t = t.AddDate(0, 0, 1)
			}
		}
	}()

	return out
}

func (d DateRange) Err() error {
	return d.err
}

var (
	ErrDateRangeEmpty            = errors.New("empty aoc date range")
	ErrDateRangeMissingSeperator = errors.New("missing seperator")
	ErrDateRangeDateNotAnAoC     = errors.New("date not an aoc date")
)

func ParseDateRange(s string) (DateRange, error) {
	return parseDateRange(s, nowLocal)
}
func parseDateRange(s string, nowf func() time.Time) (DateRange, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return DateRange{}, ErrDateRangeEmpty
	}

	if s == ":" {
		return DateRange{
			Start: time.Date(FirstYear, time.December, int(FirstDay), 0, 0, 0, 0, Timezone),
			End:   latestAoCTime(nowf),
			err:   nil,
		}, nil
	}

	if !strings.ContainsRune(s, ':') {
		return DateRange{}, ErrDateRangeMissingSeperator
	}

	ds := strings.Split(s, ":")
	startS, endS := ds[0], ds[1]

	dr := DateRange{}

	start, errS := parseDate(startS)
	end, errE := parseDate(endS)

	if len(startS) == 0 { // ":2015-01"
		dr.Start = time.Date(FirstYear, time.December, int(FirstDay), 0, 0, 0, 0, Timezone)
		if errE != nil {
			return DateRange{err: errE}, errE
		}
		dr.End = end
	} else if len(endS) == 0 { // "2015-01:"
		dr.End = latestAoCTime(nowf)
		if errS != nil {
			return DateRange{err: errS}, errS
		}
		dr.Start = start
	} else { // "2015-01:2016:01"
		if errS != nil {
			return DateRange{err: errS}, errS
		}
		if errE != nil {
			return DateRange{err: errE}, errE
		}

		if end.Before(start) {
			dr.Start, dr.End = end, start
		} else {
			dr.Start, dr.End = start, end
		}
	}

	return dr, nil
}

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return time.Time{}, err
	}

	if t.Year() < FirstYear || t.Day() > int(LastDay) {
		return time.Time{}, ErrDateRangeDateNotAnAoC
	}

	t = time.Date(t.Year(), time.December, t.Day(), 0, 0, 0, 0, Timezone)
	return t, nil
}

func LatestAoCTime() time.Time {
	return latestAoCTime(nowLocal)
}
func latestAoCTime(nowf func() time.Time) time.Time {
	now := nowf()

	if now.Month() < time.December {
		return time.Date(now.Year()-1, time.December, int(LastDay), 0, 0, 0, 0, Timezone)
	}
	if now.Month() == time.December {
		d := Day(now.Day())
		if d > LastDay {
			d = LastDay
		}

		return time.Date(now.Year(), time.December, int(d), 0, 0, 0, 0, Timezone)
	}

	return time.Date(now.Year(), time.December, int(LastDay), 0, 0, 0, 0, Timezone)
}

func nowLocal() time.Time {
	return time.Now()
}
