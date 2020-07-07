package main

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/namespaces"
)

// This command is used to generate markdown documentation for each commands (custom or generated) of the CLI
func main() {
	commands := namespaces.GetCommands()
	for _, c := range commands.GetAll() {
		if c.Namespace == "instance" {
			fmt.Printf("%s %s %s | %s\n", c.Namespace, c.Resource, c.Verb, c.Short)
		}
	}
}
