package aoc

import (
	"strconv"
	"time"
)

type Year int

func (y Year) String() string {
	return strconv.Itoa(int(y))
}

const FirstYear = Year(2015)

// Timezone is when Eric Wastl unlocks puzzles and starts messing up the
// sleeping schedules of Europeans.
var Timezone = time.FixedZone("EST/UTC-5", -5*60*60)

func nowLocal() time.Time {
	return time.Now()
}

// Years returns all released AoCs years in ascending order. It should be
// noted that this uses the local time of the user's machine. Which when messed
// with could lead to an incorrect list of AoC years.
func Years() []Year {
	return years(nowLocal)
}

func years(nowf func() time.Time) []Year {
	now := nowf()
	currYear := Year(now.Year())

	if currYear < FirstYear {
		return []Year{}
	}

	years := make([]Year, 0, currYear-FirstYear+1)
	for y := FirstYear; y < currYear; y++ {
		years = append(years, y)
	}

	nextAoc := time.Date(int(currYear), time.December, 1, 0, 0, 0, 0, Timezone)
	if now.After(nextAoc) {
		years = append(years, currYear)
	}

	return years
}

// LastYear yields the most recent Advent of Code year.
func LastYear() Year {
	return lastYear(nowLocal)
}

func lastYear(nowf func() time.Time) Year {
	years := years(nowf)
	return years[len(years)-1]
}
