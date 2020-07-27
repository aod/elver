package aoc

import (
	"time"
)

const FirstYear = 2015

// Timezone is when Eric Wastl unlocks puzzles and starts messing up the
// sleeping schedules of Europeans.
var Timezone = time.FixedZone("EST/UTC-5", -5*60*60)

// Years returns all released AoCs years in ascending order. It should be
// noted that this uses the local time of the user's machine. Which when messed
// with could lead to an incorrect list of AoC years.
func Years() []int {
	return years(nowLocal)
}

func years(nowf func() time.Time) []int {
	now := nowf()
	currYear := now.Year()

	if currYear < FirstYear {
		return []int{}
	}

	years := make([]int, 0, currYear-FirstYear+1)
	for y := FirstYear; y < currYear; y++ {
		years = append(years, y)
	}

	nextAoc := time.Date(currYear, time.December, 1, 0, 0, 0, 0, Timezone)
	if now.After(nextAoc) {
		years = append(years, currYear)
	}

	return years
}
