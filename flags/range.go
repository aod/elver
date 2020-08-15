package flags

import (
	"fmt"
	"strconv"
)

// IntRange is used to specify an *inclusive* range for the std "flag" package.
// The result is stored in Value.
type IntRange struct {
	Value, Min, Max int
}

func (ir *IntRange) String() string {
	return strconv.Itoa(ir.Value)
}

// Set satasfies part of the flag.Value interface.
// It returns an error if Atoi of v failes or num resides outside of the
// inclusive int range.
// Otherwise the result is stored in Value.
func (ir *IntRange) Set(v string) error {
	num, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	if num < ir.Min || num > ir.Max {
		return fmt.Errorf("value not in range of (%d..%d): %d", ir.Min, ir.Max, num)
	}

	ir.Value = num
	return nil
}
