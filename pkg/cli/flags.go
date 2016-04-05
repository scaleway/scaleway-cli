// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "fmt"

// CommandListOpts holds a list of parameters
type CommandListOpts struct {
	Values *[]string
}

// NewListOpts create an empty CommandListOpts
func NewListOpts() CommandListOpts {
	var values []string
	return CommandListOpts{
		Values: &values,
	}
}

// String returns a string representation of a CommandListOpts object
func (opts *CommandListOpts) String() string {
	return fmt.Sprintf("%v", (*opts.Values))
}

// Set appends a new value to a CommandListOpts
func (opts *CommandListOpts) Set(value string) error {
	(*opts.Values) = append((*opts.Values), value)
	return nil
}
