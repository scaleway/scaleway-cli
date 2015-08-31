package pricing

import (
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
	Convey("Testing Basket.Add", t, func() {
		basket := NewBasket()
		So(basket, ShouldNotBeNil)
		So(basket.Length(), ShouldEqual, 0)

		err := basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 1))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 1)

		err = basket.Add(NewUsageByPathWithQuantity("/compute/c1/run", 42))
		So(err, ShouldBeNil)
		So(basket.Length(), ShouldEqual, 2)
	})
}
