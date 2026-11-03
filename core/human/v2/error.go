package human

import "fmt"

// UnknownFieldError is returned when a selected field path does not exist in
// the data, so callers like the printer can suggest the valid fields.
type UnknownFieldError struct {
	FieldName   string
	ValidFields []string
}

func (e *UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field %s, valid fields are: %v", e.FieldName, e.ValidFields)
}
