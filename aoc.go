package main

import (
	"errors"
	"fmt"
	"time"
)

func validYear(year int) error {
	if year < 2015 {
		return errors.New("advent of code first started on 2015")
	}

	now := time.Now()
	if year > now.Year() {
		return fmt.Errorf("the latest advent of code is %d", now.Year()-1)
	}

	if year == now.Year() {
		nextAoc := time.Date(year, time.December, 1, 0, 0, 0, 0, time.FixedZone("EST/UTC-5", -5*60*60))

		if now.Before(nextAoc) {
			diff := nextAoc.Sub(now)
			return fmt.Errorf("advent of code %d starts in %s", year, diff)
		}
	}

	return nil
}

func aocYears() []int {
	curr := time.Now().Year()
	years := make([]int, 0, curr-2015)
	for y := 2015; y < curr; y++ {
		years = append(years, y)
	}

	if validYear(curr) != nil {
		years = append(years, curr)
	}

	return years
}

func validDay(day int) error {
	if day < 1 || day > 25 {
		return fmt.Errorf("invalid day: %d", day)
	}

	return nil
}
