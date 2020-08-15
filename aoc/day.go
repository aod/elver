package aoc

import (
	"strconv"
)

// Day represents a day in an Advent of Code year.
type Day int

func (d Day) String() string {
	return strconv.Itoa(int(d))
}

// First and last Day constants.
const (
	FirstDay Day = 1
	LastDay      = FirstDay + 24
)
