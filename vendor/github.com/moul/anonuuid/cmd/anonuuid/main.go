package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/moul/anonuuid"
	"github.com/moul/anonuuid/vendor/github.com/codegangsta/cli"
)

// main is the entrypoint
func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Manfred Touron"
	app.Email = "https://github.com/moul"
	app.Version = "1.0.0-dev"
	app.Usage = "Anonymize UUIDs outputs"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "hexspeak",
			Usage: "Generate hexspeak style fake UUIDs",
		},
		cli.BoolFlag{
			Name:  "random, r",
			Usage: "Generate random fake UUIDs",
		},
		cli.BoolFlag{
			Name:  "keep-beginning",
			Usage: "Keep first part of the UUID unchanged",
		},
		cli.BoolFlag{
			Name:  "keep-end",
			Usage: "Keep last part of the UUID unchanged",
		},
		cli.StringFlag{
			Name:  "prefix, p",
			Usage: "Prefix generated UUIDs",
		},
		cli.StringFlag{
			Name:  "suffix",
			Usage: "Suffix generated UUIDs",
		},
	}

	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	anonuuid := anonuuid.New()

	anonuuid.Hexspeak = c.Bool("hexspeak")
	anonuuid.Random = c.Bool("random")
	anonuuid.Prefix = c.String("prefix")
	anonuuid.Suffix = c.String("suffix")
	anonuuid.KeepBeginning = c.Bool("keep-beginning")
	anonuuid.KeepEnd = c.Bool("keep-end")

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(anonuuid.Sanitize(line))
	}
}
