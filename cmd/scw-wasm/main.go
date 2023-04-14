package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
)

var commands *core.Commands

func getCommands() *core.Commands {
	if commands == nil {
		commands = namespaces.GetCommands()
	}
	return commands
}

func runCommand(args []string, stdout io.Writer, stderr io.Writer) chan int {
	ret := make(chan int, 1)
	go func() {
		exitCode, _, _ := core.Bootstrap(&core.BootstrapConfig{
			Args:      args,
			Commands:  getCommands(),
			BuildInfo: &core.BuildInfo{},
			Stdout:    stdout,
			Stderr:    stderr,
			Stdin:     nil,
		})
		ret <- exitCode
	}()

	return ret
}

func wasmRun(this js.Value, args []js.Value) (any, error) {
	cliArgs := []string{"scw"}
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	for argIndex, arg := range args {
		if arg.Type() != js.TypeString {
			return nil, fmt.Errorf("invalid argument at index %d", argIndex)
		}
		cliArgs = append(cliArgs, arg.String())
	}

	exitCodeChan := runCommand(cliArgs, stdout, stderr)
	exitCode := <-exitCodeChan
	if exitCode != 0 {
		errBody := stderr.String()
		fmt.Println(errBody)
		return js.ValueOf(errBody), fmt.Errorf("exit code: %d", exitCode)
	}
	outBody := stdout.String()
	fmt.Println(outBody)
	return js.ValueOf(outBody), nil
}

func main() {
	js.Global().Set("cliRun", asyncFunc(wasmRun))
	<-make(chan struct{}, 0)
	//end := runCommand(os.Args, os.Stdout, os.Stderr)
	//log.Fatalln(<-end)
}
