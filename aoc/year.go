package aoc

import (
	"time"
)

// Timezone is when Eric Wastl unlocks puzzles and starts messing up the
// sleeping schedules of Europeans.
var Timezone = time.FixedZone("EST/UTC-5", -5*60*60)

func nowLocal() time.Time {
	return time.Now()
}

// Years returns all released AoCs years in ascending order. It should be
// noted that this uses the local time of the user's machine. Which when messed
// with could lead to an incorrect list of AoC years.
func Years() []int {
	return years(nowLocal)
}

func years(nowf func() time.Time) []int {
	const startYear = 2015
	now := nowf()
	currYear := now.Year()

	if currYear < startYear {
		return []int{}
	}

	years := make([]int, 0, currYear-startYear+1)
	for y := startYear; y < currYear; y++ {
		years = append(years, y)
	}

	nextAoc := time.Date(currYear, time.December, 1, 0, 0, 0, 0, Timezone)
	if now.After(nextAoc) {
		years = append(years, currYear)
	}

	return years
}
