/*
Run your Go Advent of Code solutions with a single command.
Write your solution and Elver will take care of the rest.

How It Works

Elver uses plugin build mode to generate a `.so` file to dynamically look up
the solutions.
These must reside in an Advent of Code folder under the main package.

A solution for a day in an Advent of Code year is represented by 2 solvers
for part A and B.
All solvers are functions which satisfy the same signature where interface{}
is the output:

	func (input string) (interface{}, error)

A solver must be exported and it's name satisfy the following regex:

	(Day)([1-9]|1[0-9]|2[0-5])(A|B)

E.g.:

	func Day1A(input string) (interface{}, error) {
	    return 42, nil
	}

Solvers are workspaced by the Advent of Code year which is also used as the
folder name.

Code Example

	package main
	import "errors"
	func Day1A(input string) (interface{}, error) {
	    return 42, nil
	}
	func Day1B(input string) (interface{}, error) {
	    return nil, errors.New("Not implemented")
	}

Running Elver in the root directory will output something like the following:

	$ elver
	AOC 2015
	Day 1 A (312ns):
	42
	Day 1 B (956ns):
	[ERROR] Not implemented
*/
package main

import (
	"os"

	"github.com/aod/elver/internal/cmd"
)

func main() {
	cmd.Execute(os.Args)
}
