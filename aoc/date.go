package aoc

// Date represents an Advent of Code date
type Date struct {
	Year
	Day
}

// DatePart represents an Advent of Code date for a solution.
type DatePart struct {
	Date
	Part
}
