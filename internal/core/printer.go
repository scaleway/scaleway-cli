package core

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/human"
)

// Type defines an formatter format.
type PrinterType string

// String returns the formatter format converted in a string.
func (o *PrinterType) String() string {
	return string(*o)
}

// Set sets the formatter format from a string.
func (o *PrinterType) Set(v string) error {
	switch v {
	case string(PrinterTypeJSON):
		*o = PrinterTypeJSON
	case string(PrinterTypeHuman):
		*o = PrinterTypeJSON
	default:
		return fmt.Errorf("invalid format: %s", v)
	}

	return nil
}

// Type returns the FormatterType string type.
func (o *PrinterType) Type() string {
	return "string"
}

var (
	// PrinterTypeJSON defines a JSON formatter.
	PrinterTypeJSON = PrinterType("json")

	// PrinterTypeHuman defines a human readable formatted formatter.
	PrinterTypeHuman = PrinterType("human")
)

type PrinterConfig struct {
	Type   PrinterType
	Stdout io.Writer
	Stderr io.Writer
}

// NewPrinter returns an initialized formatter corresponding to a given FormatterType.
func NewPrinter(config *PrinterConfig) (*Printer, error) {
	printerType := config.Type
	if printerType == "" {
		printerType = PrinterTypeHuman
	}

	//if printerType != PrinterTypeJSON && printerType != PrinterTypeHuman {
	//	return nil, fmt.Errorf("invalid format: %s", config.Type)
	//}

	return &Printer{
		printerType: printerType,
		stdout:      config.Stdout,
		stderr:      config.Stderr,
	}, nil
}

type Printer struct {
	printerType PrinterType
	stdout      io.Writer
	stderr      io.Writer
}

func (p *Printer) Print(data interface{}, opt *human.MarshalOpt) error {
	if rawResult, isRawResult := data.(RawResult); isRawResult {
		_, err := p.stdout.Write(rawResult)
		return err
	}
	switch p.printerType {
	case PrinterTypeHuman:
		str, err := human.Marshal(data, opt)
		if err != nil {
			return err
		}
		if str == "" {
			return nil
		}
		if _, isError := data.(error); isError {
			_, err = fmt.Fprintln(p.stderr, str)
		} else {
			_, err = fmt.Fprintln(p.stdout, str)
		}
		return err
	case PrinterTypeJSON:
		_, implementMarshaler := data.(json.Marshaler)
		err, isError := data.(error)
		switch {
		case isError && !implementMarshaler:
			data = struct {
				Error string `json:"error"`
			}{Error: err.Error()}
		}

		if isError {
			return json.NewEncoder(p.stderr).Encode(data)
		}

		if reflect.TypeOf(data).Kind() == reflect.Slice && reflect.ValueOf(data).IsNil() {
			_, err := p.stdout.Write([]byte("[]\n"))
			return err
		}
		return json.NewEncoder(p.stdout).Encode(data)

	default:
		return fmt.Errorf("invalid format: %s", p.printerType)
	}
}
