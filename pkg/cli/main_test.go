package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/scaleway/scaleway-cli/pkg/commands"
	. "github.com/smartystreets/goconvey/convey"
)

func testHelpOutput(out string, err string) {
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
	/*
		for _, line := range strings.Split(out, "\n") {
			So(line, shouldFitInTerminal)
		}
	*/

	// FIXME: count amount of options/commands, and panic if amount is different

	// Testing stderr
	So(err, ShouldEqual, "")
}

func testHelpCommandOutput(command string, out string, err string) {
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

	// :lipstick:
	/*
		for _, line := range strings.Split(out, "\n") {
			So(line, shouldFitInTerminal)
		}
	*/

}

func TestHelp(t *testing.T) {
	Convey("Testing golang' `start(\"help\", ...)`", t, func() {
		Convey("start(\"help\")", func() {
			stdout := bytes.Buffer{}
			stderr := bytes.Buffer{}
			streams := commands.Streams{
				Stdin:  os.Stdin,
				Stdout: &stdout,
				Stderr: &stderr,
			}
			ec, err := Start([]string{"help"}, &streams)

			So(ec, ShouldEqual, 0)
			So(err, ShouldBeNil)
			testHelpOutput(stdout.String(), stderr.String())
		})

		cmds := append(publicCommands, secretCommands...)
		for _, command := range cmds {
			// FIXME: test 'start(COMMAND, "--help")'
			if command == "help" {
				continue
			}

			Convey(fmt.Sprintf("start(\"help\", \"%s\")", command), func() {
				stdout := bytes.Buffer{}
				stderr := bytes.Buffer{}
				streams := commands.Streams{
					Stdin:  os.Stdin,
					Stdout: &stdout,
					Stderr: &stderr,
				}
				ec, err := Start([]string{"help", command}, &streams)

				So(ec, ShouldEqual, 1)
				So(err, ShouldBeNil)
				testHelpCommandOutput(command, stdout.String(), stderr.String())

				// secondary help usage
				// FIXME: should check for 'scw login' first
				/*
					if command != "help" {
						// FIXME: this should works for "help" too
						secondaryStdout := bytes.Buffer{}
						secondaryStderr := bytes.Buffer{}
						secondaryStreams := commands.Streams{
							Stdin:  os.Stdin,
							Stdout: &secondaryStdout,
							Stderr: &secondaryStderr,
						}
						secondEc, secondErr := Start([]string{command, "--help"}, &secondaryStreams)
						So(ec, ShouldEqual, secondEc)
						//So(outStdout, ShouldEqual, secondOut)
						So(fmt.Sprintf("%v", err), ShouldEqual, fmt.Sprintf("%v", secondErr))
					}
				*/

			})
		}
	})

	Convey("Testing shell' `scw help`", t, func() {
		Convey("scw help", func() {
			cmd := exec.Command(scwcli, "help")
			out, ec, err := runCommandWithOutput(cmd)
			stderr := "" // FIXME: get real stderr

			// exit code
			So(ec, ShouldEqual, 0)
			So(err, ShouldBeNil)

			// streams
			testHelpOutput(out, stderr)
		})

		cmds := append(publicCommands, secretCommands...)
		for _, command := range cmds {
			// FIXME: test 'scw COMMAND --help'

			Convey(fmt.Sprintf("scw help %s", command), func() {
				cmd := exec.Command(scwcli, "help", command)
				out, ec, err := runCommandWithOutput(cmd)
				stderr := "" // FIXME: get real stderr

				// exit code
				So(ec, ShouldEqual, 1)
				So(fmt.Sprintf("%s", err), ShouldEqual, "exit status 1")

				// streams
				testHelpCommandOutput(command, out, stderr)

				// secondary help usage
				// FIXME: should check for 'scw login' first
				/*
					if command != "help" {
						// FIXME: this should works for "help" too
						secondCmd := exec.Command(scwcli, command, "--help")
						secondOut, secondEc, secondErr := runCommandWithOutput(secondCmd)
						So(out, ShouldEqual, secondOut)
						So(ec, ShouldEqual, secondEc)
						So(fmt.Sprintf("%v", err), ShouldEqual, fmt.Sprintf("%v", secondErr))
					}
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
