package args

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//
// Marshal & Unmarshal Errors
//

type DataMustBeAPointerError struct {
}

func (e *DataMustBeAPointerError) Error() string {
	return fmt.Sprintf("data must be a pointer to a struct")
}

type UnknownKindError struct {
	Kind reflect.Kind
}

func (e *UnknownKindError) Error() string {
	return fmt.Sprintf("unknown kind '%s'", e.Kind)
}

//
// Marshal Errors
//

type DataMustBeAmArshalableValueError struct {
}

func (e *DataMustBeAmArshalableValueError) Error() string {
	return fmt.Sprintf("data must be a marshable value (a scalar type or a Marshaler)")
}

type ValueIsNotMarshalableError struct {
	Interface interface{}
}

func (e *ValueIsNotMarshalableError) Error() string {
	return fmt.Sprintf("%T is not marshalable", e.Interface)
}

//
// Unmarshal Errors
//

type InvalidArgumentError struct {
	ArgumentName string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument '%s': must only contain lowercase letters, numbers or dashes", e.ArgumentName)
}

type UnknowArgumentError struct {
	ArgumentName string
}

func (e *UnknowArgumentError) Error() string {
	return fmt.Sprintf("unknown argument '%s'", e.ArgumentName)
}

type DuplicateArgumentError struct {
	ArgumentName string
}

func (e *DuplicateArgumentError) Error() string {
	return fmt.Sprintf("duplicate argument '%s'", e.ArgumentName)
}

type DataIsNilError struct {
}

func (e *DataIsNilError) Error() string {
	return fmt.Sprintf("data must be not be nil")
}

type DataIsNotAPointerError struct {
}

func (e *DataIsNotAPointerError) Error() string {
	return fmt.Sprintf("data must be a pointer")
}

type CannotSetNestedFieldError struct {
	ArgumentName string
	Interface    interface{}
}

func (e *CannotSetNestedFieldError) Error() string {
	return fmt.Sprintf("cannot set nested field %s for unmarshalable type %T", e.ArgumentName, e.Interface)
}

type MissingIndexOnArrayError struct {
}

func (e *MissingIndexOnArrayError) Error() string {
	return fmt.Sprintf("missing index on the array")
}

type InvalidIndexError struct {
	Index string
}

func (e *InvalidIndexError) Error() string {
	return fmt.Sprintf("invalid index: '%s' is not a positive integer", e.Index)
}

type MissingIndicesInArrayError struct {
	IndexToInsert int
	CurrentLength int
}

func (e *MissingIndicesInArrayError) Error() string {
	switch {
	case e.IndexToInsert-e.CurrentLength == 1:
		return fmt.Sprintf("missing index %d: all indices prior to %d must be set as well", e.CurrentLength, e.IndexToInsert)
	default:
		return fmt.Sprintf("missing indices %s: all indices prior to %d must be set as well", missingIndices(int(e.IndexToInsert), e.CurrentLength), e.IndexToInsert)
	}
}

type NoSubKeyForMapError struct {
	Value string
}

func (e *NoSubKeyForMapError) Error() string {
	return fmt.Sprintf("cannot handle map with no subkey, value '%v'", e.Value)
}

type MissingFieldNameForStructError struct {
	Interface interface{}
}

func (e *MissingFieldNameForStructError) Error() string {
	return fmt.Sprintf("cannot unmarshal a struct %T with not field name", e.Interface)
}

type CannotUnmarshalTypeError struct {
	Interface interface{}
}

func (e *CannotUnmarshalTypeError) Error() string {
	return fmt.Sprintf("don't know how to unmarshal type %T", e.Interface)
}

type InvalidValueError struct {
	Value string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("invalid value '%s': valid values are true or false", e.Value)
}

type CannotUnmarshalError struct {
	Interface interface{}
}

func (e *CannotUnmarshalError) Error() string {
	return fmt.Sprintf("%T is not unmarshalable", e.Interface)
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
