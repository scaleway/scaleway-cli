package args

import (
	"fmt"
	"strconv"
	"strings"
)

//
// Marshal & Unmarshal Errors
//

type DataMustBeAPointerError struct{}

func (e *DataMustBeAPointerError) Error() string {
	return "data must be a pointer to a struct"
}

//
// Marshal Errors
//

type DataMustBeAMarshalableValueError struct{}

func (e *DataMustBeAMarshalableValueError) Error() string {
	return "data must be a marshalable value (a scalar type or a Marshaler)"
}

type ValueIsNotMarshalableError struct {
	Interface interface{}
}

func (e *ValueIsNotMarshalableError) Error() string {
	return fmt.Sprintf("%T is not marshalable", e.Interface)
}

type NotMarshalableTypeError struct {
	Src interface{}
}

func (e *NotMarshalableTypeError) Error() string {
	return fmt.Sprintf("cannot marshal type '%T'", e.Src)
}

//
// Unmarshal Errors
//

type UnmarshalArgError struct {
	// ArgName is the name of the argument which causes the error.
	ArgName string

	// ArgValue is the value of the argument which causes the error.
	ArgValue string

	// Err is the wrapped error.
	Err error
}

func (e *UnmarshalArgError) Error() string {
	arg := e.ArgName
	if e.ArgValue != "" {
		arg += "=" + e.ArgValue
	}

	return fmt.Sprintf("cannot unmarshal arg '%s': %s", arg, e.Err)
}

func (e *UnmarshalArgError) Unwrap() error {
	return e.Err
}

type InvalidArgNameError struct{}

func (e *InvalidArgNameError) Error() string {
	return "arg name must only contain lowercase letters, numbers or dashes"
}

type UnknownArgError struct{}

func (e *UnknownArgError) Error() string {
	return "unknown argument"
}

type DuplicateArgError struct{}

func (e *DuplicateArgError) Error() string {
	return "duplicate argument"
}

type CannotSetNestedFieldError struct {
	Dest interface{}
}

func (e *CannotSetNestedFieldError) Error() string {
	return fmt.Sprintf("cannot set nested field for unmarshalable type %T", e.Dest)
}

type MissingIndexOnArrayError struct{}

func (e *MissingIndexOnArrayError) Error() string {
	return "missing index on the array"
}

type InvalidIndexError struct {
	Index string
}

func (e *InvalidIndexError) Error() string {
	return fmt.Sprintf("invalid index '%s' is not a positive integer", e.Index)
}

type MissingIndicesInArrayError struct {
	IndexToInsert int
	CurrentLength int
}

func (e *MissingIndicesInArrayError) Error() string {
	switch e.IndexToInsert - e.CurrentLength {
	case 1:
		return fmt.Sprintf(
			"missing index %d, all indices prior to %d must be set as well",
			e.CurrentLength,
			e.IndexToInsert,
		)
	default:
		return fmt.Sprintf(
			"missing indices, %s all indices prior to %d must be set as well",
			missingIndices(e.IndexToInsert, e.CurrentLength),
			e.IndexToInsert,
		)
	}
}

type MissingMapKeyError struct{}

func (e *MissingMapKeyError) Error() string {
	return "missing map key"
}

type MissingStructFieldError struct {
	Dest interface{}
}

func (e *MissingStructFieldError) Error() string {
	return fmt.Sprintf("missing field name for type %T", e.Dest)
}

type UnmarshalableTypeError struct {
	Dest interface{}
}

func (e *UnmarshalableTypeError) Error() string {
	return fmt.Sprintf("do not know how to unmarshal type %T", e.Dest)
}

type CannotUnmarshalError struct {
	Dest interface{}
	Err  error
}

func (e *CannotUnmarshalError) Error() string {
	return fmt.Sprintf("%T is not unmarshalable: %s", e.Dest, e.Err)
}

// ConflictArgError is return when two args that are in conflict with each other are used together.
// e.g cluster=prod cluster.name=test are conflicting args
type ConflictArgError struct {
	ArgName1 string
	ArgName2 string
}

func (e *ConflictArgError) Error() string {
	return fmt.Sprintf(
		"arguments '%s' and '%s' cannot be used simultaneously",
		e.ArgName1,
		e.ArgName2,
	)
}

// missingIndices returns a string of all the missing indices between index and length.
// e.g.: missingIndices(index=5, length=0) should return "0,1,2,3"
// e.g.: missingIndices(index=5, length=2) should return "2,3"
// e.g.: missingIndices(index=99999, length=0) should return "0,1,2,3,4,5,6,7,8,9,..."
func missingIndices(index, length int) string {
	s := []string(nil)
	for i := length; i < index; i++ {
		if i-length == 10 {
			s = append(s, "...")

			break
		}
		s = append(s, strconv.Itoa(i))
	}

	return strings.Join(s, ",")
}

type CannotParseDateError struct {
	ArgValue               string
	AbsoluteTimeParseError error
	RelativeTimeParseError error
}

func (e *CannotParseDateError) Error() string {
	return "date parsing error: could not parse " + e.ArgValue
}

type CannotParseBoolError struct {
	Value string
}

func (e *CannotParseBoolError) Error() string {
	return "invalid boolean value"
}
