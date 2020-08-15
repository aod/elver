package aoc

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Year represents an Advent of Code year.
type Year int

func (y Year) String() string {
	return strconv.Itoa(int(y))
}

// FindDir tries to find the Advent of Code year dir given the cwd.
func (y Year) FindDir(cwd string) (string, error) {
	path := filepath.Join(cwd, y.String())
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}

// FirstYear is the year when the first Advent of Code was released.
const FirstYear = Year(2015)

// AdventYears represents a collection of Advent of Code years
type AdventYears []Year

// FirstYearDir TODO
func (years AdventYears) FirstYearDir(cwd string) (Year, string, error) {
	for _, y := range years {
		p, err := y.FindDir(cwd)
		if err == nil {
			return y, p, nil
		}
	}
	return 0, "", errors.New("no aoc year dir found")
}

func (years AdventYears) Len() int           { return len(years) }
func (years AdventYears) Swap(i, j int)      { years[i], years[j] = years[j], years[i] }
func (years AdventYears) Less(i, j int) bool { return years[i] < years[j] }

// Timezone is when Eric Wastl unlocks the Advent of Code puzzles.
var Timezone = time.FixedZone("EST/UTC-5", -5*60*60)

func nowLocal() time.Time {
	return time.Now()
}

// Years returns all released AoCs years in ascending order. It should be
// noted that this uses the local time of the user's machine. Which when messed
// with could lead to an incorrect list of AoC years.
func Years() AdventYears {
	return years(nowLocal)
}

func years(nowf func() time.Time) AdventYears {
	now := nowf()
	currYear := Year(now.Year())

	if currYear < FirstYear {
		return AdventYears{}
	}

	years := make(AdventYears, 0, currYear-FirstYear+1)
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
