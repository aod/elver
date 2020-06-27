package flags

import (
	"fmt"
	"strconv"
)

type IntRange struct {
	Value, Min, Max int
}

func (ir *IntRange) String() string {
	return strconv.Itoa(ir.Value)
}

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
