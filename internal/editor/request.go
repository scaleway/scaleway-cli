package editor

import "reflect"

// createGetRequest creates a GetRequest from given type and populate it with content from updateRequest
func createGetRequest(updateRequest interface{}, getRequestType reflect.Type) interface{} {
	updateRequestV := reflect.ValueOf(updateRequest)

	getRequest := reflect.New(getRequestType).Interface()
	getRequestV := reflect.ValueOf(getRequest)

	// Fill GetRequest args using Update arg content
	// This should copy important argument like ID, zone
	valueMapper(getRequestV, updateRequestV)

	return getRequest
}

// copyAndCompleteUpdateRequest return a copy of updateRequest completed with resource content
func copyAndCompleteUpdateRequest(updateRequest interface{}, resource interface{}) interface{} {
	resourceV := reflect.ValueOf(resource)
	updateRequestV := reflect.ValueOf(updateRequest)

	// Create a new updateRequest that will be edited
	// It will allow user to edit it, then we will extract diff to perform update
	newUpdateRequestV := reflect.New(updateRequestV.Type().Elem())
	valueMapper(newUpdateRequestV, updateRequestV)
	valueMapper(newUpdateRequestV, resourceV)

	return newUpdateRequestV.Interface()
}

func newRequest(request interface{}) interface{} {
	requestType := reflect.TypeOf(request)

	if requestType.Kind() == reflect.Pointer {
		requestType = requestType.Elem()
	}

	return reflect.New(requestType).Interface()
}
