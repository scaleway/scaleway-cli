package pricing

import (
	"math/big"
	"testing"
	"time"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
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
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(72, 100)) // 0.72

			err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(42, 1)))
			So(err, ShouldBeNil)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(144, 100)) // 1.44

			err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", big.NewRat(600, 1)))
			So(err, ShouldBeNil)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(744, 100)) // 7.44
		})
		Convey("1 compute instance with 2 volumes and 1 ip", func() {
			basket := NewBasket()

			basket.Add(NewUsageByPath("/compute/c1/run"))
			basket.Add(NewUsageByPath("/ip/dynamic"))
			basket.Add(NewUsageByPath("/storage/local/ssd/storage"))
			basket.Add(NewUsageByPath("/storage/local/ssd/storage"))
			So(basket.Length(), ShouldEqual, 4)

			basket.SetDuration(1 * time.Minute)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(136, 100)) // 1.36

			basket.SetDuration(1 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(136, 100)) // 1.36

			basket.SetDuration(2 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(232, 100)) // 2.32

			basket.SetDuration(2 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(8399, 1000)) // 8.399

			basket.SetDuration(30 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(11999, 1000)) // 11.999

			// FIXME: this test is false, the capacity is per month
			basket.SetDuration(365 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, big.NewRat(11999, 1000)) // 11.999
		})
	})
}
