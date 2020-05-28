package printer

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/human"
)

// jsonPrinter is the JSON implementation of Formatter.
type jsonPrinter struct {
	Writer      io.Writer
	ErrorWriter io.Writer
}

// Print prints data in JSON format.
func (o *jsonPrinter) Print(data interface{}, opt *human.MarshalOpt) error {
	rawBody, isStandardError := standardErrorJSON(data)
	_, implementMarshaler := data.(json.Marshaler)
	err, isError := data.(error)
	switch {
	case isStandardError:
		data = rawBody
	case isError && !implementMarshaler:
		data = struct {
			Error string `json:"error"`
		}{Error: err.Error()}
	}

	if isError || isStandardError {
		err := json.NewEncoder(o.ErrorWriter).Encode(data)
		if err != nil {
			return err
		}
	}

	if reflect.TypeOf(data).Kind() == reflect.Slice && data == nil {
		_, err := o.Writer.Write([]byte("[]"))
		return err
	}
	return json.NewEncoder(o.Writer).Encode(data)
}
