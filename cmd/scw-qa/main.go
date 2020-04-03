package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/namespaces"
	"github.com/scaleway/scaleway-cli/internal/qa"
	"github.com/scaleway/scaleway-cli/internal/tabwriter"
	"github.com/scaleway/scaleway-cli/internal/terminal"
)

func main() {
	commands := namespaces.GetCommands()
	errors := qa.LintCommands(commands)

	fmt.Println(terminal.Style("Errors:", color.Bold))
	for _, err := range errors {
		fmt.Printf("%v\n", err)
	}

	errorCounts := map[string]int{}
	for _, err := range errors {
		errorCounts[fmt.Sprintf("%T", err)]++
	}

	fmt.Println(terminal.Style("\nSummary:", color.Bold))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for key, count := range errorCounts {
		_, _ = fmt.Fprintf(w, "%s\t%d\n", strings.TrimPrefix(key, "*qa."), count)
	}
	_ = w.Flush()
}
