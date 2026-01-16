package alias

import "strings"

type Alias struct {
	// alias' key
	Name string
	// whole command
	Command []string

	// list of args in the command
	args []string
}

func (a *Alias) Args() []string {
	if a.args == nil {
		a.computeArgs()
	}

	return a.args
}

func (a *Alias) computeArgs() {
	a.args = []string{}
	for _, cmd := range a.Command {
		before, _, ok := strings.Cut(cmd, "=")
		if ok {
			a.args = append(a.args, before)
		}
	}
}
