package main

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/command"
	"github.com/scaleway/scaleway-cli/internal/qa"
)

func main() {
	commands := command.GetCommands()
	errors := qa.LintCommands(commands)

	errorCounts := map[string]int{}
	for _, err := range errors {
		switch t := err.(type) {
		default:
			errorCounts[fmt.Sprintf("%T", t)] += 1
		}
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
