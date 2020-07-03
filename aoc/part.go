package aoc

type Part rune

func (p Part) String() string {
	return string(p)
}

const (
	Part1 Part = iota + 'A'
	Part2
)
