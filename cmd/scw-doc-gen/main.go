package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scaleway/scaleway-cli/internal/docgen"
	"github.com/scaleway/scaleway-cli/internal/namespaces"
)

func main() {
	commands := namespaces.GetCommands()

	outDir := flag.String("outdir", "./docs/commands", "Directory where markdown will be created")
	flag.Parse()

	stats, err := os.Stat(*outDir)
	if err != nil {
		panic(err)
	}

	if !stats.IsDir() {
		panic(fmt.Errorf("outdir must be a valid directory"))
	}

	err = docgen.GenerateDocs(commands, *outDir)
	if err != nil {
		panic(err)
	}
}
