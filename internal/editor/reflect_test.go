package editor

import (
	"net"
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_valueMapper(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg1 string
		Arg2 int
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.Equal(t, src.Arg1, dest.Arg1)
	assert.Equal(t, src.Arg2, dest.Arg2)
}

func Test_valueMapperInvalidType(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg1 string
		Arg2 bool
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.Equal(t, src.Arg1, dest.Arg1)
	assert.Zero(t, dest.Arg2)
}

func Test_valueMapperDifferentFields(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg3 string
		Arg4 bool
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.Zero(t, dest.Arg3)
	assert.Zero(t, dest.Arg4)
}

func Test_valueMapperPointers(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg1 *string
		Arg2 *int
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.NotNil(t, dest.Arg1)
	assert.EqualValues(t, src.Arg1, *dest.Arg1)
	assert.NotNil(t, dest.Arg2)
	assert.EqualValues(t, src.Arg2, *dest.Arg2)
}

func Test_valueMapperPointersWithPointers(t *testing.T) {
	src := struct {
		Arg1 *string
		Arg2 *int32
	}{scw.StringPtr("1"), scw.Int32Ptr(1)}
	dest := struct {
		Arg1 *string
		Arg2 *int32
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.NotNil(t, dest.Arg1)
	assert.EqualValues(t, src.Arg1, dest.Arg1)
	assert.NotNil(t, dest.Arg2)
	assert.EqualValues(t, src.Arg2, dest.Arg2)
}

func Test_valueMapperSlice(t *testing.T) {
	src := struct {
		Arg1 []string
		Arg2 []int
	}{
		[]string{"1", "2", "3"},
		[]int{1, 2, 3},
	}
	dest := struct {
		Arg1 []string
		Arg2 []int
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.NotNil(t, dest.Arg1)
	assert.EqualValues(t, src.Arg1, dest.Arg1)
	assert.NotNil(t, dest.Arg2)
	assert.EqualValues(t, src.Arg2, dest.Arg2)
}

func Test_valueMapperSliceOfPointers(t *testing.T) {
	src := struct {
		Arg1 []string
		Arg2 []int
	}{
		[]string{"1", "2", "3"},
		[]int{1, 2, 3},
	}
	dest := struct {
		Arg1 []*string
		Arg2 []*int
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.NotNil(t, dest.Arg1)
	assert.Equal(t, len(src.Arg1), len(dest.Arg1))
	for i := range src.Arg1 {
		assert.NotNil(t, dest.Arg1[i])
		assert.Equal(t, src.Arg1[i], *dest.Arg1[i])
	}

	assert.NotNil(t, dest.Arg2)
	assert.Equal(t, len(src.Arg2), len(dest.Arg2))
	for i := range src.Arg2 {
		assert.NotNil(t, dest.Arg2[i])
		assert.Equal(t, src.Arg2[i], *dest.Arg2[i])
	}
}

func Test_valueMapperSliceStructPointer(t *testing.T) {
	_, ipnet, err := net.ParseCIDR("192.168.0.0/24")
	assert.Nil(t, err)

	src := instance.ListSecurityGroupRulesResponse{
		TotalCount: 0,
		Rules: []*instance.SecurityGroupRule{
			{
				ID:        "id",
				Protocol:  "protocol",
				Direction: "direction",
				Action:    "action",
				IPRange: scw.IPNet{
					IPNet: *ipnet,
				},
				DestPortFrom: scw.Uint32Ptr(1000),
				DestPortTo:   scw.Uint32Ptr(2000),
				Position:     12,
				Editable:     true,
				Zone:         "zone",
			},
		},
	}
	dest := instance.SetSecurityGroupRulesRequest{
		Rules: nil,
	}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.NotNil(t, dest.Rules)
	assert.Equal(t, 1, len(dest.Rules))
	expectedRule := src.Rules[0]
	actualRule := dest.Rules[0]
	assert.NotNil(t, actualRule.ID)
	assert.Equal(t, expectedRule.ID, *actualRule.ID)
	assert.Equal(t, expectedRule.Protocol, actualRule.Protocol)
	assert.Equal(t, expectedRule.Direction, actualRule.Direction)
	assert.Equal(t, expectedRule.Action, actualRule.Action)
	assert.Equal(t, expectedRule.IPRange, actualRule.IPRange)
	assert.NotNil(t, actualRule.DestPortFrom)
	assert.Equal(t, expectedRule.DestPortFrom, actualRule.DestPortFrom)
	assert.NotNil(t, actualRule.DestPortTo)
	assert.Equal(t, expectedRule.DestPortTo, actualRule.DestPortTo)
	assert.Equal(t, expectedRule.Position, actualRule.Position)
	assert.NotNil(t, actualRule.Editable)
	assert.Equal(t, expectedRule.Editable, *actualRule.Editable)
	assert.NotNil(t, actualRule.Zone)
	assert.Equal(t, expectedRule.Zone, *actualRule.Zone)
}
