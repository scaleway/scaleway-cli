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

		err := basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 1))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 1)

		err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 42))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 2)

		err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 600))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 3)
	})
}

func TestBasket_Total(t *testing.T) {
	Convey("Testing Basket.Total", t, FailureContinues, func() {
		Convey("3 compute instances", func() {
			basket := NewBasket()
			So(basket, ShouldNotBeNil)
			current := ratZero
			difference := ratZero
			So(basket.Total(), ShouldEqualBigRat, current)

			err := basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 1))
			So(err, ShouldBeNil)
			difference = new(big.Rat).Mul(big.NewRat(60, 1), new(big.Rat).SetFloat64(0.012))
			current = new(big.Rat).Add(current, difference)
			So(basket.Total(), ShouldEqualBigRat, current)

			err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 42))
			So(err, ShouldBeNil)
			difference = new(big.Rat).Mul(big.NewRat(60, 1), new(big.Rat).SetFloat64(0.012))
			//difference = big.NewRat(72, 100)
			current = new(big.Rat).Add(current, difference)
			So(basket.Total(), ShouldEqualBigRat, current)

			err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 600))
			So(err, ShouldBeNil)
			difference = big.NewRat(6, 1)
			current = new(big.Rat).Add(current, difference)
			So(basket.Total(), ShouldEqualBigRat, current)
		})
		Convey("1 compute instance with 2 volumes and 1 ip", func() {
			basket := NewBasket()

			basket.Add(NewUsageByPath("/compute/c1/run"))
			basket.Add(NewUsageByPath("/ip/dynamic"))
			basket.Add(NewUsageByPath("/storage/local/ssd/storage"))
			basket.Add(NewUsageByPath("/storage/local/ssd/storage"))
			So(basket.Length(), ShouldEqual, 4)

			basket.SetDuration(1 * time.Minute)
			So(basket.Total(), ShouldEqualBigRat, new(big.Rat).SetFloat64(1.36))

			basket.SetDuration(1 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, new(big.Rat).SetFloat64(1.36))

			basket.SetDuration(2 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, new(big.Rat).SetFloat64(2.32))

			basket.SetDuration(2 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, new(big.Rat).SetFloat64(8.39))

			basket.SetDuration(30 * 24 * time.Hour)
			So(basket.Total(), ShouldEqualBigRat, new(big.Rat).SetFloat64(1.99))
		})
	})
}
