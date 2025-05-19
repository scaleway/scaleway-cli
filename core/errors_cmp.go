package core

import (
	"errors"
	"net/http"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

func IsNotFoundError(err error) bool {
	notFoundError := &scw.ResourceNotFoundError{}
	responseError := &scw.ResponseError{}

	return errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound ||
		errors.As(err, &notFoundError)
}
