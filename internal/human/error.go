package human

import "fmt"

type UnknownFieldError struct {
	FieldName   string
	ValidFields []string
}

func (u *UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field %s, valid fields are: %s", u.FieldName, u.ValidFields)
}
