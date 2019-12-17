package args

import (
	"fmt"
)

func duplicateArgumentError(argName string) error {
	return fmt.Errorf("duplicate argument '%v='", argName)
}
