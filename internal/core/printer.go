package core

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/human"
)

// Type defines an formatter format.
type PrinterType string

func (p PrinterType) String() string {
	return string(p)
}

const (
	// PrinterTypeJSON defines a JSON formatter.
	PrinterTypeJSON = PrinterType("json")

	// PrinterTypeHuman defines a human readable formatted formatter.
	PrinterTypeHuman = PrinterType("human")

	// Option to enable pretty output on json printer.
	PrinterOptJSONPretty = "pretty"
)

type PrinterConfig struct {
	OutputFlag string
	Stdout     io.Writer
	Stderr     io.Writer
}

// NewPrinter returns an initialized formatter corresponding to a given FormatterType.
func NewPrinter(config *PrinterConfig) (*Printer, error) {
	printer := &Printer{
		stdout: config.Stdout,
		stderr: config.Stderr,
	}

	// First we parse OutputFlag to extract printerName and printerOpt (e.g json=pretty)
	tmp := strings.SplitN(config.OutputFlag, "=", 2)
	printerName := tmp[0]
	printerOpt := ""
	if len(tmp) > 1 {
		printerOpt = tmp[1]
	}

	// We call the correct setup method depending on the printer type
	switch printerName {
	case PrinterTypeHuman.String():
		setupHumanPrinter(printer, printerOpt)
	case PrinterTypeJSON.String():
		err := setupJSONPrinter(printer, printerOpt)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid output format: %s", printerName)
	}

	return printer, nil
}

func setupJSONPrinter(printer *Printer, opts string) error {
	printer.printerType = PrinterTypeJSON
	switch opts {
	case PrinterOptJSONPretty:
		printer.jsonPretty = true
	case "":
	default:
		return fmt.Errorf("invalid option %s for json outout. Valid options are: %s", opts, PrinterOptJSONPretty)
	}
	return nil
}

func setupHumanPrinter(printer *Printer, opts string) {
	printer.printerType = PrinterTypeHuman
	if opts != "" {
		printer.humanFields = strings.Split(opts, ",")
	}
}

type Printer struct {
	printerType PrinterType
	stdout      io.Writer
	stderr      io.Writer

	// Enable pretty print on json output
	jsonPretty bool

	// Allow to select specifics column in a table with human printer
	humanFields []string
}

func (p *Printer) Print(data interface{}, opt *human.MarshalOpt) error {
	// No matter the printer type if data is a RawResult we should print it as is.
	if rawResult, isRawResult := data.(RawResult); isRawResult {
		_, err := p.stdout.Write(rawResult)
		return err
	}

	var err error
	switch p.printerType {
	case PrinterTypeHuman:
		err = p.printHuman(data, opt)
	case PrinterTypeJSON:
		err = p.printJSON(data)
	default:
		err = fmt.Errorf("unknown format: %s", p.printerType)
	}

	if err != nil {
		// Only try to print error using the printer if data is not already an error to avoid infinite recursion
		if _, isError := data.(error); !isError {
			return p.Print(err, nil)
		}
		return err
	}
	return nil
}

func (p *Printer) printHuman(data interface{}, opt *human.MarshalOpt) error {
	_, isError := data.(error)

	if !isError {
		if opt == nil {
			opt = &human.MarshalOpt{}
		}

		if len(p.humanFields) > 0 && reflect.TypeOf(data).Kind() != reflect.Slice {
			return fmt.Errorf("list of fields for human output is only supported for commands that return a list")
		}

		if len(p.humanFields) > 0 {
			opt.Fields = []*human.MarshalFieldOpt(nil)
			for _, field := range p.humanFields {
				opt.Fields = append(opt.Fields, &human.MarshalFieldOpt{
					FieldName: field,
				})
			}
		}
	}

	str, err := human.Marshal(data, opt)
	switch e := err.(type) {
	case *human.UnknownFieldError:
		return &CliError{
			Err:  fmt.Errorf("unknown field '%s' in output options", e.FieldName),
			Hint: fmt.Sprintf("Valid fields are: %s", strings.Join(e.ValidFields, ", ")),
		}
	case nil:
		// Do nothing
	default:
		return err
	}

	// If human marshal return an empty string we avoid printing empty line
	if str == "" {
		return nil
	}

	if _, isError := data.(error); isError {
		_, err = fmt.Fprintln(p.stderr, str)
	} else {
		_, err = fmt.Fprintln(p.stdout, str)
	}
	return err
}

func (p *Printer) printJSON(data interface{}) error {
	_, implementMarshaler := data.(json.Marshaler)
	err, isError := data.(error)

	if isError && !implementMarshaler {
		data = map[string]string{
			"error": err.Error(),
		}
	}

	writer := p.stdout
	if isError {
		writer = p.stderr
	}
	encoder := json.NewEncoder(writer)
	if p.jsonPretty {
		encoder.SetIndent("", "  ")
	}

	// We handle special case to make sure that a nil slice is marshal as `[]`
	if reflect.TypeOf(data).Kind() == reflect.Slice && reflect.ValueOf(data).IsNil() {
		_, err := p.stdout.Write([]byte("[]\n"))
		return err
	}

	return encoder.Encode(data)
}
