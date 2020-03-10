package main

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/namespaces"

	"github.com/scaleway/scaleway-cli/internal/qa"
)

func main() {
	commands := namespaces.GetCommands()
	errors := qa.LintCommands(commands)

	errorCounts := map[string]int{}
	for _, err := range errors {
		errorCounts[fmt.Sprintf("%T", err)]++
	}

	fmt.Printf("Errors:\n")
	for _, err := range errors {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("\nSummary:\n")
	for key, count := range errorCounts {
		fmt.Printf("%s: %d\n", key, count)
	}
}
