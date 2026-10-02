package main

import (
	"bytes"
	"context"
	"reflect"
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
			args = populateSliceFields(args)
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

// populateSliceFields initializes slice fields with a single zero-value element
// to allow template execution with index access (e.g., {{ index .ServerIDs 0 }})
func populateSliceFields(args any) any {
	v := reflect.ValueOf(args)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return args
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Slice && field.CanSet() {
			elemType := field.Type().Elem()
			// Create a new element (zero value)
			newElem := reflect.New(elemType).Elem()
			// Set the slice to contain one element
			field.Set(reflect.Append(field, newElem))
		}
	}
	return args
}
