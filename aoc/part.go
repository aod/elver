package aoc

// Part represents a part in a single Advent of Code day.
type Part rune

func (p Part) String() string {
	return string(p)
}

// First and last Part constants.
const (
	Part1 Part = iota + 'A'
	Part2
)
