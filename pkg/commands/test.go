package commands

import (
	"fmt"
	"strings"

	"github.com/pborman/uuid"
)

func shouldBeAnUUID(actual interface{}, expected ...interface{}) string {
	input := actual.(string)
	input = strings.TrimSpace(input)
	uuid := uuid.Parse(input)
	if uuid == nil {
		return fmt.Sprintf("%q should be an UUID", actual)
	}
	return ""
}
