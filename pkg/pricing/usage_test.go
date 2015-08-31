package pricing

import (
	"fmt"
	"math/big"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ShouldEqualBigRat(actual interface{}, expected ...interface{}) string {
	actualRat := actual.(*big.Rat)
	expectRat := expected[0].(*big.Rat)
	cmp := actualRat.Cmp(expectRat)
	if cmp == 0 {
		return ""
	}

	output := fmt.Sprintf("big.Rat are not matching: %q != %q\n", actualRat, expectRat)

	actualFloat64, _ := actualRat.Float64()
	expectFloat64, _ := expectRat.Float64()
	output += fmt.Sprintf("                          %f != %f", actualFloat64, expectFloat64)
	return output
}

func TestNewUsageByPath(t *testing.T) {
	Convey("Testing NewUsageByPath()", t, func() {
		usage := NewUsageByPath("/compute/c1/run", 1)
		So(usage.PricingObject.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
	})
}

func TestNewUsage(t *testing.T) {
	Convey("Testing NewUsage()", t, func() {
		object := CurrentPricing.GetByPath("/compute/c1/run")
		usage := NewUsage(object, 1)
		So(usage.PricingObject.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
	})
}

func TestUsage_BillableQuantity(t *testing.T) {
	Convey("Testing Usage.BillableQuantity()", t, FailureContinues, func() {
		object := &PricingObject{
			UnitQuantity: 60,
		}
		usage := NewUsage(object, -1)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, -1000)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, 0)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, 1)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsage(object, 59)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsage(object, 59.9999)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsage(object, 60)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsage(object, 60.00001)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*2, 1))

		usage = NewUsage(object, 61)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*2, 1))

		usage = NewUsage(object, 119)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*2, 1))

		usage = NewUsage(object, 121)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*3, 1))

		usage = NewUsage(object, 1000)
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*17, 1))
	})
}

func TestUsage_LostQuantity(t *testing.T) {
	Convey("Testing Usage.LostQuantity()", t, FailureContinues, func() {
		object := &PricingObject{
			UnitQuantity: 60,
		}
		usage := NewUsage(object, -1)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, -1000)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, 0)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, 1)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60-1, 1))

		usage = NewUsage(object, 59)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60-59, 1))

		usage = NewUsage(object, 59.9999)
		So(usage.LostQuantity(), ShouldEqualBigRat, new(big.Rat).Sub(big.NewRat(60, 1), new(big.Rat).SetFloat64(59.9999)))

		usage = NewUsage(object, 60)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(object, 60.00001)
		So(usage.LostQuantity(), ShouldEqualBigRat, new(big.Rat).SetFloat64(60*2-60.00001))

		usage = NewUsage(object, 61)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*2-61, 1))

		usage = NewUsage(object, 119)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*2-119, 1))

		usage = NewUsage(object, 121)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*3-121, 1))

		usage = NewUsage(object, 1000)
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*17-1000, 1))
	})
}

func TestUsage_Total(t *testing.T) {
	Convey("Testing Usage.Total()", t, FailureContinues, func() {
		object := PricingObject{
			UnitQuantity: 60,
			UnitPrice:    0.012,
			UnitPriceCap: 6,
		}

		usage := NewUsage(&object, -1)
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(&object, 0)
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsage(&object, 1)
		So(usage.Total(), ShouldEqualBigRat, new(big.Rat).Mul(big.NewRat(60, 1), new(big.Rat).SetFloat64(0.012)))

		usage = NewUsage(&object, 61)
		So(usage.Total(), ShouldEqualBigRat, new(big.Rat).Mul(big.NewRat(120, 1), new(big.Rat).SetFloat64(0.012)))

		usage = NewUsage(&object, 1000)
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(6, 1))
	})
}
