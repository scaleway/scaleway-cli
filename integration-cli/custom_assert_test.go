package integrationcli

import "fmt"

func shouldFitInTerminal(actual interface{}, expected ...interface{}) string {
	if len(actual.(string)) < 80 {
		return ""
	}
	return fmt.Sprintf("len(%q)\n -> %d chars (> 80 chars)", actual, len(actual.(string)))
}
