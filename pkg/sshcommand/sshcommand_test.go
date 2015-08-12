package sshcommand

import "fmt"

func ExampleCommand() *Command {
	return &Command{
		Host: "1.2.3.4",
	}
}

func ExampleCommand_options() *Command {
	return &Command{
		SkipHostKeyChecking: true,
		Host:                "1.2.3.4",
		Quiet:               true,
		AllocateTTY:         true,
		Command:             []string{"echo", "hello world"},
		Debug:               true,
	}
}

func ExampleCommand_complex() *Command {
	return &Command{
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
}

func ExampleCommand_gateway() *Command {
	return &Command{
		Host:    "1.2.3.4",
		Gateway: New("5.6.7.8"),
	}
}

func ExampleCommand_New() *Command {
	return New("1.2.3.4")
}

func ExampleCommand_String() {
	fmt.Println(New("1.2.3.4").String())
	// Output: ssh 1.2.3.4
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
	// ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no 1.2.3.4 -t -t -- /bin/sh -e -x -c "\"\\\"echo\\\" \\\"hello world\\\"\""
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
	// ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o "ProxyCommand=ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -W %h:%p -l toor 5.6.7.8 -t -t" 1.2.3.4 -t -t -- /bin/sh -e -c "\"\\\"echo\\\" \\\"hello world\\\"\""
}

func ExampleCommand_Slice() {
	fmt.Println(New("1.2.3.4").Slice())
	// Output: [ssh 1.2.3.4]
}

func ExampleCommand_Slice_user() {
	fmt.Println(New("root@1.2.3.4").Slice())
	// Output: [ssh -l root 1.2.3.4]
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
	// ["ssh" "-q" "-o" "UserKnownHostsFile=/dev/null" "-o" "StrictHostKeyChecking=no" "-l" "root" "1.2.3.4" "-t" "-t" "--" "/bin/sh" "-e" "-x" "-c" "\"\\\"echo\\\" \\\"hello world\\\"\""]
}

func ExampleCommand_Slice_gateway() {
	command := Command{
		Host:    "1.2.3.4",
		Gateway: New("5.6.7.8"),
	}
	fmt.Printf("%q\n", command.Slice())

	// Output:
	// ["ssh" "-o" "ProxyCommand=ssh -W %h:%p 5.6.7.8" "1.2.3.4"]
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
	// ["ssh" "-q" "-o" "UserKnownHostsFile=/dev/null" "-o" "StrictHostKeyChecking=no" "-o" "ProxyCommand=ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -W %h:%p -l toor 5.6.7.8 -t -t" "1.2.3.4" "-t" "-t" "--" "/bin/sh" "-e" "-c" "\"echo hello world\""]
}
