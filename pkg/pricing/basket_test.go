package pricing

import (
	"math/big"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBasket(t *testing.T) {
	Convey("Testing NewBasket()", t, func() {
		basket := NewBasket()
		So(basket, ShouldNotBeNil)
		So(basket.Length(), ShouldEqual, 0)
	})
}

func TestBasket_Add(t *testing.T) {
	Convey("Testing Basket.Add", t, FailureContinues, func() {
		basket := NewBasket()
		So(basket, ShouldNotBeNil)
		So(basket.Length(), ShouldEqual, 0)

		err := basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(1, 1)))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 1)

		err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(42, 1)))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 2)

		err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(600, 1)))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 3)
	})
}

func TestBasket_Total(t *testing.T) {
	Convey("Testing Basket.Total", t, FailureContinues, func() {
		Convey("3 compute instances", func() {
			basket := NewBasket()
			So(basket, ShouldNotBeNil)
			So(basket.Total(), ShouldEqualBigRat, ratZero)

			err := basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(1, 1)))
			So(err, ShouldBeNil)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(2, 1000)) // 0.002

			err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(42, 1)))
			So(err, ShouldBeNil)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(4, 1000)) // 0.004

			err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(600, 1)))
			So(err, ShouldBeNil)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(24, 1000)) // 0.024
		})
		Convey("1 server of each type", func() {
			basket := NewBasket()
			So(basket, ShouldNotBeNil)
			So(basket.Total(), ShouldEqualBigRat, ratZero)

			basket.Add(NewUsageByPath("/compute/c1/run"))
			basket.Add(NewUsageByPath("/compute/c2s/run"))
			basket.Add(NewUsageByPath("/compute/c2m/run"))
			basket.Add(NewUsageByPath("/compute/c2l/run"))
			basket.Add(NewUsageByPath("/compute/vc1s/run"))
			basket.Add(NewUsageByPath("/compute/vc1m/run"))
			basket.Add(NewUsageByPath("/compute/vc1l/run"))

			basket.SetDuration(1 * time.Minute)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(116, 1000)) // 0.116

			basket.SetDuration(1 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(116, 1000)) // 0.116

			basket.SetDuration(2 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(232, 1000)) // 0.232

			basket.SetDuration(24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(2784, 1000)) // 2.784

			basket.SetDuration(30 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(58000, 1000)) // 58

			// FIXME: this test if false, the capacity is per month
			basket.SetDuration(365 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(58000, 1000)) // 58
		})
		Convey("1 compute instance with 2 volumes and 1 ip", func() {
			basket := NewBasket()

			basket.Add(NewUsageByPath("/compute/c1/run"))
			basket.Add(NewUsageByPath("/ip/dynamic"))
			basket.Add(NewUsageByPath("/storage/local/ssd/storage"))
			basket.Add(NewUsageByPath("/storage/local/ssd/storage"))
			So(basket.Length(), ShouldEqual, 4)

			basket.SetDuration(1 * time.Minute)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(8, 1000)) // 0.008

			basket.SetDuration(1 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(8, 1000)) // 0.008

			basket.SetDuration(2 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(12, 1000)) // 0.012

			basket.SetDuration(2 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(196, 1000)) // 0.196

			basket.SetDuration(30 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(2050, 1000)) // 2.05

			// FIXME: this test is false, the capacity is per month
			basket.SetDuration(365 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(2694, 1000)) // 2.694
		})
	})
}
