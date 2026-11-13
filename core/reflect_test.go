package core_test

import (
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

type RequestEmbedding struct {
	EmbeddingField1 string
	EmbeddingField2 int
}

type CreateRequest struct {
	*RequestEmbedding
	CreateField1 string
	CreateField2 int
}

type ExtendedRequest struct {
	*CreateRequest
	ExtendedField1 string
	ExtendedField2 int
}

type ArrowRequest struct {
	PrivateNetwork *PrivateNetwork
}

type SpecialRequest struct {
	*RequestEmbedding
	TabRequest []*ArrowRequest
}

type EndpointSpecPrivateNetwork struct {
	PrivateNetworkID string
	ServiceIP        *scw.IPNet
}

type PrivateNetwork struct {
	*EndpointSpecPrivateNetwork
	OtherValue string
}

func Test_getValuesForFieldByName(t *testing.T) {
	type TestCase struct {
		cmdArgs        any
		fieldName      string
		expectedError  string
		expectedValues []reflect.Value
	}

	expectedServiceIP := &scw.IPNet{
		IP:   net.ParseIP("192.0.2.1"),
		Mask: net.CIDRMask(24, 32),
	}

	tests := []struct {
		name     string
		testCase TestCase
		testFunc func(*testing.T, TestCase)
	}{
		{
			name: "Simple test",
			testCase: TestCase{
				cmdArgs: &ExtendedRequest{
					CreateRequest: &CreateRequest{
						RequestEmbedding: &RequestEmbedding{
							EmbeddingField1: "value1",
							EmbeddingField2: 2,
						},
						CreateField1: "value3",
						CreateField2: 4,
					},
					ExtendedField1: "value5",
					ExtendedField2: 6,
				},
				fieldName:      "EmbeddingField1",
				expectedError:  "",
				expectedValues: []reflect.Value{reflect.ValueOf("value1")},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				values, err := core.GetValuesForFieldByName(
					reflect.ValueOf(tc.cmdArgs),
					strings.Split(tc.fieldName, "."),
				)
				if err != nil {
					assert.Equal(t, tc.expectedError, err.Error())
				} else if tc.expectedValues != nil && !reflect.DeepEqual(tc.expectedValues[0].Interface(), values[0].Interface()) {
					t.Errorf(
						"Expected %v, got %v",
						tc.expectedValues[0].Interface(),
						values[0].Interface(),
					)
				}
			},
		},
		{
			name: "Error test",
			testCase: TestCase{
				cmdArgs: &ExtendedRequest{
					CreateRequest: &CreateRequest{
						RequestEmbedding: &RequestEmbedding{
							EmbeddingField1: "value1",
							EmbeddingField2: 2,
						},
						CreateField1: "value3",
						CreateField2: 4,
					},
					ExtendedField1: "value5",
					ExtendedField2: 6,
				},
				fieldName:      "NotExist",
				expectedError:  "field NotExist does not exist for ExtendedRequest",
				expectedValues: []reflect.Value{reflect.ValueOf("value1")},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				values, err := core.GetValuesForFieldByName(
					reflect.ValueOf(tc.cmdArgs),
					strings.Split(tc.fieldName, "."),
				)
				if err != nil {
					assert.Equal(t, tc.expectedError, err.Error())
				} else if tc.expectedValues != nil && !reflect.DeepEqual(tc.expectedValues[0].Interface(), values[0].Interface()) {
					t.Errorf(
						"Expected %v, got %v",
						tc.expectedValues[0].Interface(),
						values[0].Interface(),
					)
				}
			},
		},
		{
			name: "Special test",
			testCase: TestCase{
				cmdArgs: &SpecialRequest{
					RequestEmbedding: &RequestEmbedding{
						EmbeddingField1: "value1",
						EmbeddingField2: 2,
					},
					TabRequest: []*ArrowRequest{
						{
							PrivateNetwork: &PrivateNetwork{
								EndpointSpecPrivateNetwork: &EndpointSpecPrivateNetwork{
									ServiceIP: &scw.IPNet{
										IP:   net.ParseIP("192.0.2.1"),
										Mask: net.CIDRMask(24, 32),
									},
								},
							},
						},
						{
							PrivateNetwork: &PrivateNetwork{
								OtherValue: "hello",
							},
						},
					},
				},
				fieldName:      "tabRequest.{index}.privateNetwork.serviceIP",
				expectedError:  "",
				expectedValues: []reflect.Value{reflect.ValueOf(expectedServiceIP)},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				values, err := core.GetValuesForFieldByName(
					reflect.ValueOf(tc.cmdArgs),
					strings.Split(tc.fieldName, "."),
				)
				if err != nil {
					assert.Nil(t, err.Error())
				} else if tc.expectedValues != nil && !reflect.DeepEqual(tc.expectedValues[0].Interface(), values[0].Interface()) {
					t.Errorf(
						"Expected %v, got %v",
						tc.expectedValues[0].Interface(),
						values[0].Interface(),
					)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.testCase)
		})
	}
}

func Test_ResponseTypeFields(t *testing.T) {
	type server struct {
		ID   string
		Name string
		Tags []string
	}

	tests := map[string]struct {
		responseType reflect.Type
		expected     []string // nil when the type is not struct-based
	}{
		"struct":       {reflect.TypeFor[server](), []string{"ID", "Name", "Tags"}},
		"pointer":      {reflect.TypeFor[*server](), []string{"ID", "Name", "Tags"}},
		"slice":        {reflect.TypeFor[[]server](), []string{"ID", "Name", "Tags"}},
		"slicePointer": {reflect.TypeFor[[]*server](), []string{"ID", "Name", "Tags"}},
		"string":       {reflect.TypeFor[string](), nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fields := core.ResponseTypeFields(tc.responseType)
			if tc.expected == nil {
				assert.Nil(t, fields)
				return
			}
			names := make([]string, 0, len(fields))
			for _, f := range fields {
				names = append(names, f.Name)
			}
			assert.Equal(t, tc.expected, names)
		})
	}
}
