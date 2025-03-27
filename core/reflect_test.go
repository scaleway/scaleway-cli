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
		cmdArgs        interface{}
		fieldName      string
		expectedError  string
		expectedValues []reflect.Value
	}

	expectedServiceIP := &scw.IPNet{
		IPNet: net.IPNet{
			IP:   net.ParseIP("192.0.2.1"),
			Mask: net.CIDRMask(24, 32),
		},
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
					t.Errorf("Expected %v, got %v", tc.expectedValues[0].Interface(), values[0].Interface())
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
					t.Errorf("Expected %v, got %v", tc.expectedValues[0].Interface(), values[0].Interface())
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
										IPNet: net.IPNet{
											IP:   net.ParseIP("192.0.2.1"),
											Mask: net.CIDRMask(24, 32),
										},
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
					t.Errorf("Expected %v, got %v", tc.expectedValues[0].Interface(), values[0].Interface())
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
