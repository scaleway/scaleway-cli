package core

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/scaleway/scaleway-cli/v2/internal/gofields"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-cli/v2/internal/terraform"
)

// Type defines an formatter format.
type PrinterType string

func (p PrinterType) String() string {
	return string(p)
}

const (
	// PrinterTypeJSON defines a JSON formatter.
	PrinterTypeJSON = PrinterType("json")

	// PrinterTypeYAML defines a YAML formatter.
	PrinterTypeYAML = PrinterType("yaml")

	// PrinterTypeYAML defines a Terraform formatter.
	PrinterTypeTerraform = PrinterType("terraform")

	// PrinterTypeHuman defines a human readable formatted formatter.
	PrinterTypeHuman = PrinterType("human")

	// PrinterTypeWide defines a human-readable formatted formatter without shrinking.
	PrinterTypeWide = PrinterType("wide")

	// PrinterTypeTemplate defines a go template to use to format output.
	PrinterTypeTemplate = PrinterType("template")

	// Option to enable pretty output on json printer.
	PrinterOptJSONPretty = "pretty"

	// Option to enable pretty output on json printer.
	PrinterOptTerraformWithChildren = "with-children"
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
	case PrinterTypeWide.String():
		setupWidePrinter(printer, printerOpt)
	case PrinterTypeJSON.String():
		err := setupJSONPrinter(printer, printerOpt)
		if err != nil {
			return nil, err
		}
	case PrinterTypeYAML.String():
		printer.printerType = PrinterTypeYAML
	case PrinterTypeTerraform.String():
		err := setupTerraformPrinter(printer, printerOpt)
		if err != nil {
			return nil, err
		}
	case PrinterTypeTemplate.String():
		err := setupTemplatePrinter(printer, printerOpt)
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

func setupTerraformPrinter(printer *Printer, opts string) error {
	printer.printerType = PrinterTypeTerraform
	switch opts {
	case PrinterOptTerraformWithChildren:
		printer.terraformWithChildren = true
	case "":
	default:
		return fmt.Errorf("invalid option %s for terraform outout. Valid options are: %s", opts, PrinterOptTerraformWithChildren)
	}

	terraformVersion, err := terraform.GetLocalClientVersion()
	if err != nil {
		return err
	}

	if terraformVersion.Major < 1 || (terraformVersion.Major == 1 && terraformVersion.Minor < 5) {
		return fmt.Errorf("terraform version %s is not supported. Please upgrade to terraform >= 1.5.0", terraformVersion.String())
	}

	return nil
}

func setupTemplatePrinter(printer *Printer, opts string) error {
	printer.printerType = PrinterTypeTemplate
	if opts == "" {
		return &CliError{
			Err:     fmt.Errorf("cannot use a template output with an empty template"),
			Hint:    `Try using golang template string: scw instance server list -o template="{{ .ID }} ☜(˚▽˚)☞ {{ .Name }}"`,
			Details: `https://golang.org/pkg/text/template`,
		}
	}

	t, err := template.New("OutputFormat").Parse(opts)
	if err != nil {
		return err
	}
	printer.template = t

	return nil
}

func setupHumanPrinter(printer *Printer, opts string) {
	printer.printerType = PrinterTypeHuman
	if opts != "" {
		printer.humanFields = strings.Split(opts, ",")
	}
}

func setupWidePrinter(printer *Printer, opts string) {
	setupHumanPrinter(printer, opts)
	printer.printerType = PrinterTypeWide
}

type Printer struct {
	printerType PrinterType
	stdout      io.Writer
	stderr      io.Writer

	// Enable pretty print on json output
	jsonPretty bool

	// Enable children fetching on terraform output
	terraformWithChildren bool

	// go template to use on template output
	template *template.Template

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
	case PrinterTypeWide:
		err = p.printWide(data, opt)
	case PrinterTypeJSON:
		err = p.printJSON(data)
	case PrinterTypeYAML:
		err = p.printYAML(data)
	case PrinterTypeTerraform:
		err = p.printTerraform(data)
	case PrinterTypeTemplate:
		err = p.printTemplate(data)
	default:
		err = fmt.Errorf("unknown format: %s", p.printerType)
	}

	if err != nil {
		// if the printer itself returns an error, don't try to format it just print it
		_, err := fmt.Fprintln(p.stderr, err.Error())
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
			return p.printHuman(fmt.Errorf("list of fields for human output is only supported for commands that return a list"), nil)
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
		return p.printHuman(&CliError{
			Err:  fmt.Errorf("unknown field '%s' in output options", e.FieldName),
			Hint: fmt.Sprintf("Valid fields are: %s", strings.Join(e.ValidFields, ", ")),
		}, nil)
	case nil:
		// Do nothing
	default:
		return err
	}

	// If human marshal return an empty string we avoid printing empty line
	if str == "" {
		return nil
	}

	if isError {
		_, err = fmt.Fprintln(p.stderr, str)
	} else {
		_, err = fmt.Fprintln(p.stdout, str)
	}
	return err
}

func (p *Printer) printWide(data interface{}, opt *human.MarshalOpt) error {
	if opt != nil {
		opt.DisableShrinking = true
	} else {
		opt = &human.MarshalOpt{
			DisableShrinking: true,
		}
	}
	return p.printHuman(data, opt)
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

func (p *Printer) printYAML(data interface{}) error {
	_, implementMarshaler := data.(yaml.Marshaler)
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
	encoder := yaml.NewEncoder(writer)

	return encoder.Encode(data)
}

func (p *Printer) printTerraform(data interface{}) error {
	writer := p.stdout
	if _, isError := data.(error); isError {
		return p.printHuman(data, nil)
	}

	hcl, err := terraform.GetHCL(data)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(hcl))
	return err
}

func (p *Printer) printTemplate(data interface{}) error {
	writer := p.stdout
	if _, isError := data.(error); isError {
		return p.printHuman(data, nil)
	}

	dataValue := reflect.ValueOf(data)
	switch dataValue.Type().Kind() {
	// If we have a slice of value, we apply the template for each item
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			elemValue := dataValue.Index(i)
			err := p.template.Execute(writer, elemValue)
			if err != nil {
				return p.printHuman(&CliError{
					Err:     err,
					Message: "templating error",
					Hint:    fmt.Sprintf("Acceptable values are:\n  - %s", strings.Join(gofields.ListFields(elemValue.Type()), "\n  - ")),
				}, nil)
			}
			_, err = writer.Write([]byte{'\n'})
			if err != nil {
				return err
			}
		}
	default:
		err := p.template.Execute(writer, data)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}
