package pricing

import (
	"fmt"
	"math/big"
)

func ShouldEqualBigRat(actual interface{}, expected ...interface{}) string {
	actualRat := actual.(*big.Rat)
	expectRat := expected[0].(*big.Rat)
	cmp := actualRat.Cmp(expectRat)
	if cmp == 0 {
		return ""
	}

	output := fmt.Sprintf("big.Rat are not matching: %q != %q\n", actualRat, expectRat)

	actualFloat64, _ := actualRat.Float64()
	expectFloat64, _ := expectRat.Float64()
	output += fmt.Sprintf("                          %f != %f", actualFloat64, expectFloat64)
	return output
}
