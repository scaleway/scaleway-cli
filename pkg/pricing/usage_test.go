package pricing

import (
	"math/big"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUsageByPathWithQuantity(t *testing.T) {
	Convey("Testing NewUsageByPathWithQuantity()", t, func() {
		usage := NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(1, 1))
		So(usage.Object.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
	})
}

func TestNewUsageByPath(t *testing.T) {
	Convey("Testing NewUsageByPath()", t, func() {
		usage := NewUsageByPath("/compute/c1/run")
		So(usage.Object.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqualBigRat, ratZero)
	})
}

func TestNewUsageWithQuantity(t *testing.T) {
	Convey("Testing NewUsageWithQuantity()", t, func() {
		object := CurrentPricing.GetByPath("/compute/c1/run")
		usage := NewUsageWithQuantity(object, big.NewRat(1, 1))
		So(usage.Object.Path, ShouldEqual, "/compute/c1/run")
		So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
	})
}

func TestUsage_SetStartEnd(t *testing.T) {
	Convey("Testing Usage.SetStartEnd()", t, func() {
		object := Object{
			UsageGranularity: time.Minute,
		}
		usage := NewUsage(&object)
		layout := "2006-Jan-02 15:04:05"
		start, err := time.Parse(layout, "2015-Jan-25 13:15:42")
		So(err, ShouldBeNil)
		end, err := time.Parse(layout, "2015-Jan-25 13:16:10")
		So(err, ShouldBeNil)
		err = usage.SetStartEnd(start, end)
		So(err, ShouldBeNil)
		So(usage.Quantity, ShouldEqualBigRat, big.NewRat(2, 1))
	})
}

func TestUsage_SetDuration(t *testing.T) {
	Convey("Testing Usage.SetDuration()", t, FailureContinues, func() {
		Convey("UsageGranularity=time.Minute", func() {
			object := Object{
				UsageGranularity: time.Minute,
			}
			usage := NewUsage(&object)

			err := usage.SetDuration(time.Minute * 10)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(10, 1))

			err = usage.SetDuration(time.Minute + time.Second)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(2, 1))

			err = usage.SetDuration(0 * time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(0, 1))

			err = usage.SetDuration(-1 * time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(0, 1))

			err = usage.SetDuration(10*time.Hour + 5*time.Minute + 10*time.Second)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(60*10+5+1, 1))

			err = usage.SetDuration(10 * time.Nanosecond)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
		})

		Convey("UsageGranularity=time.Hour", func() {
			object := Object{
				UsageGranularity: time.Hour,
			}
			usage := NewUsage(&object)

			err := usage.SetDuration(time.Minute * 10)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))

			err = usage.SetDuration(time.Minute + time.Second)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))

			err = usage.SetDuration(0 * time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(0, 1))

			err = usage.SetDuration(-1 * time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(0, 1))

			err = usage.SetDuration(10*time.Hour + 5*time.Minute + 10*time.Second)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(11, 1))

			err = usage.SetDuration(10 * time.Nanosecond)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
		})

		Convey("UsageGranularity=time.Hour*24", func() {
			object := Object{
				UsageGranularity: time.Hour * 24,
			}
			usage := NewUsage(&object)

			err := usage.SetDuration(time.Minute * 10)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))

			err = usage.SetDuration(time.Minute + time.Second)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))

			err = usage.SetDuration(0 * time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(0, 1))

			err = usage.SetDuration(-1 * time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(0, 1))

			err = usage.SetDuration(10*time.Hour + 5*time.Minute + 10*time.Second)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))

			err = usage.SetDuration(3*24*time.Hour + 1*time.Minute)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(4, 1))

			err = usage.SetDuration(10 * time.Nanosecond)
			So(err, ShouldBeNil)
			So(usage.Quantity, ShouldEqualBigRat, big.NewRat(1, 1))
		})
	})
}

func TestUsage_BillableQuantity(t *testing.T) {
	Convey("Testing Usage.BillableQuantity()", t, FailureContinues, func() {
		object := &Object{
			UnitQuantity: big.NewRat(60, 1),
		}
		usage := NewUsageWithQuantity(object, big.NewRat(-1, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(-1000, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(0, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(1, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(59, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(599999, 10000)) // 59.9999
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(60, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(6000001, 100000)) // 60.00001
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*2, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(61, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*2, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(119, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*2, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(121, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*3, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(1000, 1))
		So(usage.BillableQuantity(), ShouldEqualBigRat, big.NewRat(60*17, 1))
	})
}

func TestUsage_LostQuantity(t *testing.T) {
	Convey("Testing Usage.LostQuantity()", t, FailureContinues, func() {
		object := &Object{
			UnitQuantity: big.NewRat(60, 1),
		}
		usage := NewUsageWithQuantity(object, big.NewRat(-1, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(-1000, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(0, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(1, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60-1, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(59, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60-59, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(599999, 10000)) // 59.9999
		So(usage.LostQuantity(), ShouldEqualBigRat, new(big.Rat).Sub(big.NewRat(60, 1), big.NewRat(599999, 10000)))

		usage = NewUsageWithQuantity(object, big.NewRat(60, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(6000001, 100000)) // 60.00001
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(6000000*2-6000001, 100000))

		usage = NewUsageWithQuantity(object, big.NewRat(61, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*2-61, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(119, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*2-119, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(121, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*3-121, 1))

		usage = NewUsageWithQuantity(object, big.NewRat(1000, 1))
		So(usage.LostQuantity(), ShouldEqualBigRat, big.NewRat(60*17-1000, 1))
	})
}

func TestUsage_Total(t *testing.T) {
	Convey("Testing Usage.Total()", t, FailureContinues, func() {
		object := Object{
			UnitQuantity: big.NewRat(60, 1),
			UnitPrice:    big.NewRat(12, 1000), // 0.012
			UnitPriceCap: big.NewRat(6, 1),
		}

		usage := NewUsageWithQuantity(&object, big.NewRat(-1, 1))
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(&object, big.NewRat(0, 1))
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(0, 1))

		usage = NewUsageWithQuantity(&object, big.NewRat(1, 1))
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(12, 1000)) // 0.012

		usage = NewUsageWithQuantity(&object, big.NewRat(61, 1))
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(24, 1000)) // 0.024

		usage = NewUsageWithQuantity(&object, big.NewRat(1000, 1))
		So(usage.Total(), ShouldEqualBigRat, big.NewRat(204, 1000)) // 0.204
	})
}
