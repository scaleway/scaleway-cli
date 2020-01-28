package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/human"
)

// Type defines an formatter format.
type Type string

// String returns the formatter format converted in a string.
func (o *Type) String() string {
	return string(*o)
}

// Set sets the formatter format from a string.
func (o *Type) Set(v string) error {
	*o = Type(v)
	return nil
}

// Type returns the FormatterType string type.
func (o *Type) Type() string {
	return "string"
}

var (
	// JSON defines a JSON formatter.
	JSON = Type("json")

	// Human defines a human readable formatted formatter.
	Human = Type("human")
)

// Printer contains all the formatter logic for a given format.
type Printer interface {
	Print(data interface{}, opt *human.MarshalOpt) error
}

// New returns an initialized formatter corresponding to a given FormatterType.
func New(printerType Type, writer io.Writer, errorWriter io.Writer) (Printer, error) {
	if printerType == "" {
		printerType = Human
	}

	printerType = Type(strings.ToLower(string(printerType)))

	switch printerType {
	case JSON:
		return &jsonPrinter{
			Writer:      writer,
			ErrorWriter: errorWriter,
		}, nil
	case Human:
		return &humanPrinter{
			Writer:      writer,
			ErrorWriter: errorWriter,
		}, nil
	default:
		return nil, fmt.Errorf("invalid format: %s", printerType)
	}
}
