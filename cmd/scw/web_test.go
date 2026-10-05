package main

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/commands"
)

func Test_WebValidateTemplates(t *testing.T) {
	cmds := commands.GetCommands(context.Background())

	// Test that web urls are valid templates
	type failedTemplate struct {
		Cmd string
		Err error
	}
	errs := []any(nil)

	for _, cmd := range cmds.GetSortedCommand() {
		if cmd.WebURL == "" {
			continue
		}
		_, err := template.New("").Parse(cmd.WebURL)
		if err != nil {
			errs = append(errs, failedTemplate{
				Cmd: cmd.GetCommandLine("scw"),
				Err: err,
			})
		}
	}
	if len(errs) > 0 {
		t.Fatal(errs...)
	}
}

func Test_WebValidateTemplatesVariables(t *testing.T) {
	cmds := commands.GetCommands(context.Background())

	// Test that web urls are valid templates
	type failedTemplate struct {
		Cmd string
		Err error
	}
	errs := []any(nil)

	for _, cmd := range cmds.GetSortedCommand() {
		if cmd.WebURL == "" {
			continue
		}
		tmpl, err := template.New("").Parse(cmd.WebURL)
		if err != nil {
			continue
		}
		var args any
		if cmd.ArgsType != nil {
			args = reflect.New(cmd.ArgsType).Interface()
		}

		err = tmpl.Execute(bytes.NewBuffer(nil), args)
		// Temporarily skip errors related to accessing elements from empty slices/arrays, these are valid templates that work with real data.
		// Necessary until #6342 is done.
		if err != nil && !isSliceIndexOutOfRangeError(err) {
			errs = append(errs, failedTemplate{
				Cmd: cmd.GetCommandLine("scw"),
				Err: err,
			})
		}
	}
	if len(errs) > 0 {
		t.Fatal(errs...)
	}
}

func isSliceIndexOutOfRangeError(err error) bool {
	return strings.Contains(err.Error(), "slice index out of range")
}
