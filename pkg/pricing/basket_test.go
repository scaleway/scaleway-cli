package pricing

import (
	"math/big"
	"testing"

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
		current = new(big.Rat).Add(current, difference)
		So(basket.Total(), ShouldEqualBigRat, current)

		err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 600))
		So(err, ShouldBeNil)
		difference = big.NewRat(6, 1)
		current = new(big.Rat).Add(current, difference)
		So(basket.Total(), ShouldEqualBigRat, current)
	})
}
