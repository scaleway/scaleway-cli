package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/qa"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

func main() {
	cmds := commands.GetCommands(context.Background())
	errors := qa.LintCommands(cmds)

	fmt.Println(terminal.Style("Errors:", color.Bold))
	for _, err := range errors {
		fmt.Printf("%v\n", err)
	}

	errorCounts := map[string]int{}
	for _, err := range errors {
		errorCounts[fmt.Sprintf("%T", err)]++
	}

	fmt.Println(terminal.Style("\nSummary:", color.Bold))
	grid := make([][]string, 0, len(errorCounts))
	for key, count := range errorCounts {
		grid = append(grid, []string{strings.TrimPrefix(key, "*qa."), strconv.Itoa(count)})
	}
	if len(grid) > 0 {
		fmt.Println(human.RenderGrid(grid, 0))
	}
}
