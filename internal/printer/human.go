package printer

import (
	"fmt"
	"io"

	"github.com/scaleway/scaleway-cli/internal/ui/human"
)

// humanPrinter is the human readable implementation of Printer.
type humanPrinter struct {
	Writer      io.Writer
	ErrorWriter io.Writer
}

func (o *humanPrinter) Print(data interface{}, opt *human.MarshalOpt) error {
	str, err := human.Marshal(data, opt)

	if err != nil {
		_, err = fmt.Fprintln(o.ErrorWriter, err)
		return err
	}

	if str == "" {
		return nil
	}

	if _, isError := data.(error); isError {
		_, err = fmt.Fprintln(o.ErrorWriter, str)
		return err
	}

	_, err = fmt.Fprintln(o.Writer, str)
	return err
}
