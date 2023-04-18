package main

import "os"

type Args struct {
	callback     string
	targetObject string
}

func getArgs() Args {
	args := Args{}
	if len(os.Args) > 0 {
		args.callback = os.Args[0]
	}
	if len(os.Args) > 1 {
		args.targetObject = os.Args[1]
	}

	return args
}
