package main

import (
	"bytes"
	"reflect"
	"testing"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/commands"
)

func Test_WebValidateTemplates(t *testing.T) {
	cmds := commands.GetCommands()

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
	cmds := commands.GetCommands()

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
