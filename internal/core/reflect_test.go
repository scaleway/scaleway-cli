package core_test

import (
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
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
	*PrivateNetwork
}

type SpecialRequest struct {
	*RequestEmbedding
	TabRequest []*ArrowRequest
}

type EndpointSpecPrivateNetwork struct {
	PrivateNetworkId string
	ServiceIP        *scw.IPNet
}

type PrivateNetwork struct {
	*EndpointSpecPrivateNetwork
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
			Mask: net.CIDRMask(24, 32), // Exemple pour un masque de sous-réseau /24
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
				values, err := core.GetValuesForFieldByName(reflect.ValueOf(tc.cmdArgs), strings.Split(tc.fieldName, "."))
				if err != nil {
					assert.Equal(t, tc.expectedError, err.Error())
				} else {
					if tc.expectedValues != nil && !reflect.DeepEqual(tc.expectedValues[0].Interface(), values[0].Interface()) {
						t.Errorf("Expected %v, got %v", tc.expectedValues[0].Interface(), values[0].Interface())
					}
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
				values, err := core.GetValuesForFieldByName(reflect.ValueOf(tc.cmdArgs), strings.Split(tc.fieldName, "."))
				if err != nil {
					assert.Equal(t, tc.expectedError, err.Error())
				} else {
					if tc.expectedValues != nil && !reflect.DeepEqual(tc.expectedValues[0].Interface(), values[0].Interface()) {
						t.Errorf("Expected %v, got %v", tc.expectedValues[0].Interface(), values[0].Interface())
					}
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
								EndpointSpecPrivateNetwork: &EndpointSpecPrivateNetwork{
									ServiceIP: &scw.IPNet{
										IPNet: net.IPNet{
											IP:   net.ParseIP("198.51.100.1"),
											Mask: net.CIDRMask(24, 32), // Un autre exemple avec un masque de sous-réseau /24
										},
									},
								},
							},
						},
					},
				},
				fieldName:      "tabRequest.{index}.PrivateNetwork.EndpointSpecPrivateNetwork.ServiceIP",
				expectedError:  "",
				expectedValues: []reflect.Value{reflect.ValueOf(expectedServiceIP)},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				values, err := core.GetValuesForFieldByName(reflect.ValueOf(tc.cmdArgs), strings.Split(tc.fieldName, "."))
				if err != nil {
					assert.Equal(t, nil, err.Error())
				} else {
					if tc.expectedValues != nil && !reflect.DeepEqual(tc.expectedValues[0].Interface(), values[0].Interface()) {
						t.Errorf("Expected %v, got %v", tc.expectedValues[0].Interface(), values[0].Interface())
					}
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
