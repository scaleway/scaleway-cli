package sshcommand

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleCommand() {
	cmd := Command{
		Host: "1.2.3.4",
	}

	// Do stuff
	fmt.Println(cmd)
}

func ExampleCommand_options() {
	cmd := Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		Debug:               true,
	}

	// Do stuff
	fmt.Println(cmd)
}

func ExampleCommand_complex() {
	cmd := Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		Gateway: &Command{
			Host:        "5.6.7.8",
			User:        "toor",
			Quiet:       true,
			AllocateTTY: true,
		},
	}

	// Do stuff
	fmt.Println(cmd)
}

func ExampleCommand_gateway() {
	cmd := Command{
		Host:    "1.2.3.4",
		Gateway: New("5.6.7.8"),
	}

	// Do stuff
	fmt.Println(cmd)
}

func ExampleNew() {
	cmd := New("1.2.3.4")

	// Do stuff
	fmt.Println(cmd)
}

func ExampleCommand_String() {
	fmt.Println(New("1.2.3.4").String())
	// Output: ssh 1.2.3.4 -p 22
}

func ExampleCommand_String_options() {
	command := Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		Debug:               true,
	}
	fmt.Println(command.String())

	// Output:
	// ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no 1.2.3.4 -t -t -p 22 -- /bin/sh -e -x -c "\"\\\"echo\\\" \\\"hello world\\\"\""
}

func ExampleCommand_String_complex() {
	command := Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		Gateway: &Command{
			SkipHostKeyChecking: true,
			Host:                "5.6.7.8",
			User:                "toor",
			Quiet:               true,
			AllocateTTY:         true,
		},
	}
	fmt.Println(command.String())

	// Output:
	// ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o "ProxyCommand=ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -W %h:%p -l toor 5.6.7.8 -t -t -p 22" 1.2.3.4 -t -t -p 22 -- /bin/sh -e -c "\"\\\"echo\\\" \\\"hello world\\\"\""
}

func ExampleCommand_Slice() {
	fmt.Println(New("1.2.3.4").Slice())
	// Output: [ssh 1.2.3.4 -p 22]
}

func ExampleCommand_Slice_user() {
	fmt.Println(New("root@1.2.3.4").Slice())
	// Output: [ssh -l root 1.2.3.4 -p 22]
}

func ExampleCommand_Slice_options() {
	command := Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		Debug:               true,
		User:                "root",
	}
	fmt.Printf("%q\n", command.Slice())

	// Output:
	// ["ssh" "-q" "-o" "UserKnownHostsFile=/dev/null" "-o" "StrictHostKeyChecking=no" "-l" "root" "1.2.3.4" "-t" "-t" "-p" "22" "--" "/bin/sh" "-e" "-x" "-c" "\"\\\"echo\\\" \\\"hello world\\\"\""]
}

func ExampleCommand_Slice_gateway() {
	command := Command{
		Host:    "1.2.3.4",
		Gateway: New("5.6.7.8"),
	}
	fmt.Printf("%q\n", command.Slice())

	// Output:
	// ["ssh" "-o" "ProxyCommand=ssh -W %h:%p 5.6.7.8 -p 22" "1.2.3.4" "-p" "22"]
}

func ExampleCommand_Slice_complex() {
	command := Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		NoEscapeCommand:     true,
		Gateway: &Command{
			SkipHostKeyChecking: true,
			Host:                "5.6.7.8",
			User:                "toor",
			Quiet:               true,
			AllocateTTY:         true,
		},
	}
	fmt.Printf("%q\n", command.Slice())

	// Output:
	// ["ssh" "-q" "-o" "UserKnownHostsFile=/dev/null" "-o" "StrictHostKeyChecking=no" "-o" "ProxyCommand=ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -W %h:%p -l toor 5.6.7.8 -t -t -p 22" "1.2.3.4" "-t" "-t" "-p" "22" "--" "/bin/sh" "-e" "-c" "\"echo hello world\""]
}

func TestCommand_defaults(t *testing.T) {
	Convey("Testing Command{} default values", t, func() {
		command := Command{}
		So(command.Host, ShouldEqual, "")
		So(command.Port, ShouldEqual, 0)
		So(command.User, ShouldEqual, "")
		So(command.SkipHostKeyChecking, ShouldEqual, false)
		So(command.Quiet, ShouldEqual, false)
		So(len(command.SSHOptions), ShouldEqual, 0)
		So(command.Gateway, ShouldEqual, nil)
		So(command.AllocateTTY, ShouldEqual, false)
		So(len(command.Command), ShouldEqual, 0)
		So(command.Debug, ShouldEqual, false)
		So(command.NoEscapeCommand, ShouldEqual, false)
		So(command.isGateway, ShouldEqual, false)
	})
}

func TestCommand_applyDefaults(t *testing.T) {
	Convey("Testing Command.applyDefaults()", t, func() {
		Convey("On a Command{}", func() {
			command := Command{}
			command.applyDefaults()
			So(command.Host, ShouldEqual, "")
			So(command.Port, ShouldEqual, 22)
			So(command.User, ShouldEqual, "")
			So(command.SkipHostKeyChecking, ShouldEqual, false)
			So(command.Quiet, ShouldEqual, false)
			So(len(command.SSHOptions), ShouldEqual, 0)
			So(command.Gateway, ShouldEqual, nil)
			So(command.AllocateTTY, ShouldEqual, false)
			So(len(command.Command), ShouldEqual, 0)
			So(command.Debug, ShouldEqual, false)
			So(command.NoEscapeCommand, ShouldEqual, false)
			So(command.isGateway, ShouldEqual, false)
		})
		Convey("On a New(\"example.com\")", func() {
			command := New("example.com")
			command.applyDefaults()
			So(command.Host, ShouldEqual, "example.com")
			So(command.Port, ShouldEqual, 22)
			So(command.User, ShouldEqual, "")
			So(command.SkipHostKeyChecking, ShouldEqual, false)
			So(command.Quiet, ShouldEqual, false)
			So(len(command.SSHOptions), ShouldEqual, 0)
			So(command.Gateway, ShouldEqual, nil)
			So(command.AllocateTTY, ShouldEqual, false)
			So(len(command.Command), ShouldEqual, 0)
			So(command.Debug, ShouldEqual, false)
			So(command.NoEscapeCommand, ShouldEqual, false)
			So(command.isGateway, ShouldEqual, false)
		})
		Convey("On a New(\"toto@example.com\")", func() {
			command := New("toto@example.com")
			command.applyDefaults()
			So(command.Host, ShouldEqual, "example.com")
			So(command.Port, ShouldEqual, 22)
			So(command.User, ShouldEqual, "toto")
			So(command.SkipHostKeyChecking, ShouldEqual, false)
			So(command.Quiet, ShouldEqual, false)
			So(len(command.SSHOptions), ShouldEqual, 0)
			So(command.Gateway, ShouldEqual, nil)
			So(command.AllocateTTY, ShouldEqual, false)
			So(len(command.Command), ShouldEqual, 0)
			So(command.Debug, ShouldEqual, false)
			So(command.NoEscapeCommand, ShouldEqual, false)
			So(command.isGateway, ShouldEqual, false)
		})
	})
}
