package jshelpers

import (
	"fmt"
	"math"
	"syscall/js"
)

func asString(value js.Value) (string, error) {
	if value.Type() == js.TypeString {
		return value.String(), nil
	}
	return "", fmt.Errorf("value type should be string")
}

func asInt(value js.Value) (int, error) {
	if value.Type() != js.TypeNumber {
		return 0, fmt.Errorf("value type should be number")
	}
	f := value.Float()

	if f != math.Trunc(f) {
		return 0, fmt.Errorf("expected an int")
	}

	return int(f), nil
}
