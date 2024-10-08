package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/internal/docgen"
)

// This command is used to generate markdown documentation for each commands (custom or generated) of the CLI
func main() {
	cmds := commands.GetCommands()

	outDir := flag.String("outdir", "./docs/commands", "Directory where markdown will be created")
	flag.Parse()

	stats, err := os.Stat(*outDir)
	if err != nil {
		panic(err)
	}

	if !stats.IsDir() {
		panic(fmt.Errorf("outdir %s must be a valid directory", *outDir))
	}

	err = docgen.GenerateDocs(cmds, *outDir)
	if err != nil {
		panic(err)
	}
}
