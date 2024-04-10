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
											Mask: net.CIDRMask(24, 32), // Exemple pour un masque de sous-réseau /24
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

//
//func TestValidateRequiredOneOfGroups(t *testing.T) {
//	tests := []struct {
//		name          string
//		setupManager  func() *core.OneOfGroupManager
//		rawArgs       args.RawArgs
//		expectedError string
//		ArgsType      interface{}
//	}{
//		{
//			name: "Required group satisfied with first argument",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups:         map[string][]string{"group1": {"a", "b"}},
//					RequiredGroups: map[string]bool{"group1": true},
//				}
//			},
//			rawArgs: []string{"a=true"},
//			ArgsType: struct {
//				A bool
//				B bool
//			}{},
//			expectedError: "",
//		},
//		{
//			name: "Required group satisfied with second argument",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups:         map[string][]string{"group1": {"a", "b"}},
//					RequiredGroups: map[string]bool{"group1": true},
//				}
//			},
//			rawArgs: []string{"b=true"},
//			ArgsType: struct {
//				A bool
//				B bool
//			}{},
//			expectedError: "",
//		},
//		{
//			name: "Required group not satisfied",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups:         map[string][]string{"group1": {"a", "b"}},
//					RequiredGroups: map[string]bool{"group1": true},
//				}
//			},
//			rawArgs: []string{"c=true"},
//			ArgsType: struct {
//				A bool
//				B bool
//				C bool
//			}{},
//			expectedError: "at least one argument from the 'group1' group is required",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			manager := tt.setupManager()
//			err := manager.ValidateRequiredOneOfGroups(tt.rawArgs, tt.ArgsType)
//
//			if tt.expectedError == "" {
//				assert.NoError(t, err, "Expected no error, got %v", err)
//			} else {
//				assert.EqualError(t, err, tt.expectedError, fmt.Sprintf("Expected error message '%s', got '%v'", tt.expectedError, err))
//			}
//		})
//	}
//}
//
//func TestValidateUniqueOneOfGroups(t *testing.T) {
//	tests := []struct {
//		name          string
//		setupManager  func() *core.OneOfGroupManager
//		rawArgs       args.RawArgs
//		expectedError string
//		ArgsType      interface{}
//	}{
//		{
//			name: "Required group satisfied with first argument",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{"group1": {"a", "b"}},
//				}
//			},
//			rawArgs: []string{"A=true"},
//			ArgsType: struct {
//				A bool
//				B bool
//			}{},
//			expectedError: "",
//		},
//		{
//			name: "No arguments passed",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{"group1": {"a", "b"}},
//				}
//			},
//			rawArgs: []string{},
//			ArgsType: struct {
//				A bool
//				B bool
//			}{},
//			expectedError: "",
//		},
//		{
//			name: "Multiple groups, all satisfied",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{
//						"group1": {"a", "b"},
//						"group2": {"c", "d"},
//					},
//				}
//			},
//			rawArgs: []string{"a=true", "c=true"},
//			ArgsType: struct {
//				A string
//				B string
//				C string
//				D string
//			}{},
//			expectedError: "",
//		},
//		{
//			name: "Multiple groups, one satisfied",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{
//						"group1": {"a", "b"},
//						"group2": {"c", "d"},
//					},
//				}
//			},
//			rawArgs: []string{"a=true"},
//			ArgsType: struct {
//				A string
//				B string
//				C string
//				D string
//			}{},
//			expectedError: "",
//		},
//		{
//			name: "Multiple groups, not exclusive argument for groups 2",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{
//						"group1": {"a", "b"},
//						"group2": {"c", "d"},
//					},
//				}
//			},
//			rawArgs: []string{"a=true", "c=true", "d=true"},
//			ArgsType: struct {
//				A string
//				B string
//				C string
//				D string
//			}{},
//			expectedError: "arguments 'c' and 'd' are mutually exclusive",
//		},
//		{
//			name: "Multiple groups, not exclusive argument for groups 1",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{
//						"group1": {"a", "b"},
//						"group2": {"c", "d"},
//					},
//				}
//			},
//			rawArgs: []string{"a=true", "b=true", "c=true"},
//			ArgsType: struct {
//				A string
//				B string
//				C string
//				D string
//			}{},
//			expectedError: "arguments 'a' and 'b' are mutually exclusive",
//		},
//		{
//			name: "One group, not exclusive argument for groups 1",
//			setupManager: func() *core.OneOfGroupManager {
//				return &core.OneOfGroupManager{
//					Groups: map[string][]string{
//						"group1": {"a", "b", "c", "d"},
//					},
//				}
//			},
//			rawArgs: []string{"a=true", "d=true"},
//			ArgsType: struct {
//				A string
//				B string
//				C string
//				D string
//			}{},
//			expectedError: "arguments 'a' and 'd' are mutually exclusive",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			manager := tt.setupManager()
//			err := manager.ValidateUniqueOneOfGroups(tt.rawArgs, tt.ArgsType)
//			if tt.expectedError == "" {
//				assert.NoError(t, err, "Expected no error, got %v", err)
//			} else {
//				assert.EqualError(t, err, tt.expectedError, fmt.Sprintf("Expected error message '%s', got '%v'", tt.expectedError, err))
//			}
//		})
//	}
//}
//
//func Test_ValidateOneOf(t *testing.T) {
//	t.Run("Simple one-of validation check", core.Test(&core.TestConfig{
//		Commands: core.NewCommands(&core.Command{
//			Namespace:            "oneof",
//			ArgsType:             reflect.TypeOf(args.RawArgs{}),
//			AllowAnonymousClient: true,
//			Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//				return &core.SuccessResult{}, nil
//			},
//			ArgSpecs: core.ArgSpecs{
//				{
//					Name:       "a",
//					OneOfGroup: "groups1",
//				},
//				{
//					Name:       "b",
//					OneOfGroup: "groups1",
//				},
//			},
//		}),
//		Cmd: "scw oneof a=yo",
//		Check: core.TestCheckCombine(
//			core.TestCheckExitCode(0),
//		),
//	}))
//	t.Run("Required argument group check passes", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace:            "oneof",
//				ArgsType:             reflect.TypeOf(args.RawArgs{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "a",
//						OneOfGroup: "group1",
//						Required:   true,
//					},
//					{
//						Name:       "b",
//						OneOfGroup: "group1",
//						Required:   true,
//					},
//				},
//			}),
//			Cmd: "scw oneof b=yo",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(0),
//			),
//		})(t)
//	})
//
//	t.Run("Fail when required group is missing", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace:            "oneof",
//				ArgsType:             reflect.TypeOf(args.RawArgs{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "a",
//						OneOfGroup: "group1",
//						Required:   true,
//					},
//					{
//						Name:       "b",
//						OneOfGroup: "group1",
//						Required:   true,
//					},
//				},
//			}),
//			Cmd: "scw oneof c=yo",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(1),
//				core.TestCheckError(fmt.Errorf("at least one argument from the 'group1' group is required")),
//			),
//		})(t)
//	})
//
//	t.Run("Check for mutual exclusivity in arguments", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					A string
//					B string
//				}{}), AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "a",
//						OneOfGroup: "group1",
//					},
//					{
//						Name:       "b",
//						OneOfGroup: "group1",
//					},
//				},
//			}),
//			Cmd: "scw oneof a=yo b=no",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(1),
//				core.TestCheckError(fmt.Errorf("arguments 'a' and 'b' are mutually exclusive")),
//			),
//		})(t)
//	})
//
//	t.Run("Three arguments' mutual exclusivity test", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					A string
//					B string
//					C string
//				}{}), AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "a",
//						OneOfGroup: "group1",
//					},
//					{
//						Name:       "b",
//						OneOfGroup: "group1",
//					},
//					{
//						Name:       "c",
//						OneOfGroup: "group1",
//					},
//				},
//			}),
//			Cmd: "scw oneof a=yo c=no",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(1),
//				core.TestCheckError(fmt.Errorf("arguments 'a' and 'c' are mutually exclusive")),
//			),
//		})(t)
//	})
//
//	t.Run("Indexed arguments' exclusivity check", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					SSHKey     []string
//					AllSSHKeys bool
//				}{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "ssh-key.{index}",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//					{
//						Name:       "all-ssh-keys",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//				},
//			}),
//			Cmd: "scw oneof all-ssh-keys=true ssh-key.0=11111111-1111-1111-1111-111111111111",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(1),
//				core.TestCheckError(fmt.Errorf("arguments 'ssh-key.{index}' and 'all-ssh-keys' are mutually exclusive")),
//			),
//		})(t)
//	})
//
//	t.Run("Passing an indexed argument test", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					SSHKey     []string
//					AllSSHKeys bool
//				}{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "ssh-key.{index}",
//						OneOfGroup: "ssh",
//					},
//					{
//						Name:       "all-ssh-keys",
//						OneOfGroup: "ssh",
//					},
//				},
//			}),
//			Cmd: "scw oneof ssh-key.0=11111111-1111-1111-1111-111111111111",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(0),
//			),
//		})(t)
//	})
//
//	t.Run("Required indexed argument satisfies condition", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					SSHKey     []string
//					AllSSHKeys bool
//				}{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "ssh-key.{index}",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//					{
//						Name:       "all-ssh-keys",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//				},
//			}),
//			Cmd: "scw oneof ssh-key.0=11111111-1111-1111-1111-111111111111",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(0),
//			),
//		})(t)
//	})
//
//	t.Run("Exclusive all SSH keys and indexed key test", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					SSHKey     []string
//					AllSSHKeys bool
//				}{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "ssh-key.{index}",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//					{
//						Name:       "all-ssh-keys",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//				},
//			}),
//			Cmd: "scw oneof all-ssh-keys=true",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(0),
//			),
//		})(t)
//	})
//
//	t.Run("Ungrouped argument with unsatisfied group fails", func(t *testing.T) {
//		core.Test(&core.TestConfig{
//			Commands: core.NewCommands(&core.Command{
//				Namespace: "oneof",
//				ArgsType: reflect.TypeOf(struct {
//					SSHKey     []string
//					AllSSHKeys bool
//					Arg        bool
//				}{}),
//				AllowAnonymousClient: true,
//				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
//					return &core.SuccessResult{}, nil
//				},
//				ArgSpecs: core.ArgSpecs{
//					{
//						Name:       "ssh-key.{index}",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//					{
//						Name:       "all-ssh-keys",
//						OneOfGroup: "ssh",
//						Required:   true,
//					},
//					{
//						Name: "arg",
//					},
//				},
//			}),
//			Cmd: "scw oneof arg=true",
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(1),
//				core.TestCheckError(fmt.Errorf("at least one argument from the 'ssh' group is required")),
//			),
//		})(t)
//	})
//}
