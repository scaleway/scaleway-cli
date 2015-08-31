package pricing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUsageByPath(t *testing.T) {
	Convey("Testing NewUsageByPath()", t, func() {
		usage := NewUsageByPath("/compute/c1/run", 1)
		So(usage.PricingObject.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqual, 1)
	})
}

func TestNewUsage(t *testing.T) {
	Convey("Testing NewUsage()", t, func() {
		object := CurrentPricing.GetByPath("/compute/c1/run")
		usage := NewUsage(object, 1)
		So(usage.PricingObject.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqual, 1)
	})
}

func TestUsage_BillableQuantity(t *testing.T) {
	Convey("Testing Usage.BillableQuantity()", t, FailureContinues, func() {
		object := &PricingObject{
			UnitQuantity: 60,
		}
		usage := NewUsage(object, -1)
		So(usage.BillableQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, -1000)
		So(usage.BillableQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, 0)
		So(usage.BillableQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, 1)
		So(usage.BillableQuantity(), ShouldEqual, 60)

		usage = NewUsage(object, 59)
		So(usage.BillableQuantity(), ShouldEqual, 60)

		usage = NewUsage(object, 59.9999)
		So(usage.BillableQuantity(), ShouldEqual, 60)

		usage = NewUsage(object, 60)
		So(usage.BillableQuantity(), ShouldEqual, 60)

		usage = NewUsage(object, 60.00001)
		So(usage.BillableQuantity(), ShouldEqual, 60*2)

		usage = NewUsage(object, 61)
		So(usage.BillableQuantity(), ShouldEqual, 60*2)

		usage = NewUsage(object, 119)
		So(usage.BillableQuantity(), ShouldEqual, 60*2)

		usage = NewUsage(object, 121)
		So(usage.BillableQuantity(), ShouldEqual, 60*3)

		usage = NewUsage(object, 1000)
		So(usage.BillableQuantity(), ShouldEqual, 60*17)
	})
}

func TestUsage_LostQuantity(t *testing.T) {
	Convey("Testing Usage.LostQuantity()", t, FailureContinues, func() {
		object := &PricingObject{
			UnitQuantity: 60,
		}
		usage := NewUsage(object, -1)
		So(usage.LostQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, -1000)
		So(usage.LostQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, 0)
		So(usage.LostQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, 1)
		So(usage.LostQuantity(), ShouldEqual, 60-1)

		usage = NewUsage(object, 59)
		So(usage.LostQuantity(), ShouldEqual, 60-59)

		// oops error, float64 precision isn't sufficient
		// usage = NewUsage(object, 59.9999)
		// So(usage.LostQuantity(), ShouldEqual, 60-59.9999)

		usage = NewUsage(object, 60)
		So(usage.LostQuantity(), ShouldEqual, 0)

		usage = NewUsage(object, 60.00001)
		So(usage.LostQuantity(), ShouldEqual, 60*2-60.00001)

		usage = NewUsage(object, 61)
		So(usage.LostQuantity(), ShouldEqual, 60*2-61)

		usage = NewUsage(object, 119)
		So(usage.LostQuantity(), ShouldEqual, 60*2-119)

		usage = NewUsage(object, 121)
		So(usage.LostQuantity(), ShouldEqual, 60*3-121)

		usage = NewUsage(object, 1000)
		So(usage.LostQuantity(), ShouldEqual, 60*17-1000)
	})
}

func TestUsage_Total(t *testing.T) {
	Convey("Testing Usage.Total()", t, func() {
		object := PricingObject{
			UnitQuantity: 60,
			UnitPrice:    0.012,
			UnitPriceCap: 6,
		}

		usage := NewUsage(&object, -1)
		So(usage.Total(), ShouldEqual, 0)

		usage = NewUsage(&object, 0)
		So(usage.Total(), ShouldEqual, 0)

		usage = NewUsage(&object, 1)
		So(usage.Total(), ShouldEqual, 0.012*60)

		usage = NewUsage(&object, 61)
		So(usage.Total(), ShouldEqual, 0.012*120)

		usage = NewUsage(&object, 1000)
		So(usage.Total(), ShouldEqual, 6)
	})
}
