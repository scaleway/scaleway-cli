package integrationcli

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
)

func TestHelp(t *testing.T) {
	Convey("Testing 'scw help'", t, func() {
		Convey("scw help", func() {
			cmd := exec.Command(scwcli, "help")
			out, ec, err := runCommandWithOutput(cmd)
			So(ec, ShouldEqual, 0)
			So(err, ShouldBeNil)

			// headers & footers
			So(out, ShouldContainSubstring, "Usage: scw [OPTIONS] COMMAND [arg...]")
			So(out, ShouldContainSubstring, "Interact with Scaleway from the command line.")
			So(out, ShouldContainSubstring, "Run 'scw COMMAND --help' for more information on a command.")

			// options
			So(out, ShouldContainSubstring, "Options:")
			for _, option := range publicOptions {
				So(out, ShouldContainSubstring, " "+option)
			}

			// public commands
			So(out, ShouldContainSubstring, "Commands:")
			for _, command := range publicCommands {
				So(out, ShouldContainSubstring, "    "+command)
			}

			// secret commands
			for _, command := range secretCommands {
				So(out, ShouldNotContainSubstring, "    "+command)
			}

			// :lipstick:
			for _, line := range strings.Split(out, "\n") {
				if false { // FIXME: disabled for now
					So(line, shouldFitInTerminal)
				}
			}

			// FIXME: count amount of options/commands, and panic if amount is different
		})

		commands := append(publicCommands, secretCommands...)
		for _, command := range commands {
			if command == "help" {
				// FIXME: this should not be a special case
				continue
			}
			Convey(fmt.Sprintf("scw --help %s", command), func() {
				cmd := exec.Command(scwcli, command, "--help")
				_, ec, err := runCommandWithOutput(cmd)

				// exit code
				So(ec, ShouldEqual, 1)
				So(fmt.Sprintf("%s", err), ShouldEqual, "exit status 1")

				// FIXME: disabled for now because it depends of 'scw login'
				/*
					// Header
					So(out, ShouldContainSubstring, fmt.Sprintf("Usage: scw %s", command))
					// FIXME: test description
					// FIXME: test parameters

					// Options
					So(out, ShouldContainSubstring, "Options:")
					So(out, ShouldContainSubstring, " -h, --help=false")
					// FIXME: test other options

					// Examples
					// FIXME: todo
					//So(out, ShouldContainSubstring, "Examples:")

					secondCmd := exec.Command(scwcli, "help", command)
					secondOut, secondEc, secondErr := runCommandWithOutput(secondCmd)
					So(out, ShouldEqual, secondOut)
					So(ec, ShouldEqual, secondEc)
					So(fmt.Sprintf("%v", err), ShouldEqual, fmt.Sprintf("%v", secondErr))
				*/
			})
		}
		Convey("Unknown command", func() {
			cmd := exec.Command(scwcli, "boogie")
			out, ec, err := runCommandWithOutput(cmd)
			So(out, ShouldContainSubstring, "scw: unknown subcommand boogie")
			So(ec, ShouldEqual, 1)
			So(fmt.Sprintf("%s", err), ShouldEqual, "exit status 1")
		})
	})
}
