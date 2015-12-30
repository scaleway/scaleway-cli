package pricing

import (
	"math/big"
	"time"
)

// Usage represents a billing usage
type Usage struct {
	Object   *Object
	Quantity *big.Rat
}

// NewUsageByPath returns a new Usage with defaults fetched from a path
func NewUsageByPath(objectPath string) Usage {
	return NewUsageByPathWithQuantity(objectPath, ratZero)
}

// NewUsageByPathWithQuantity returns a NewUsageByPath with a specific quantity
func NewUsageByPathWithQuantity(objectPath string, quantity *big.Rat) Usage {
	return NewUsageWithQuantity(CurrentPricing.GetByPath(objectPath), quantity)
}

// NewUsageWithQuantity returns a new Usage with quantity
func NewUsageWithQuantity(object *Object, quantity *big.Rat) Usage {
	return Usage{
		Object:   object,
		Quantity: quantity,
	}
}

// NewUsage returns a new Usage
func NewUsage(object *Object) Usage {
	return NewUsageWithQuantity(object, ratZero)
}

// SetQuantity sets the quantity of an Usage
func (u *Usage) SetQuantity(quantity *big.Rat) error {
	u.Quantity = ratMax(quantity, ratZero)
	return nil
}

// SetDuration sets the duration of an Usage
func (u *Usage) SetDuration(duration time.Duration) error {
	minutes := new(big.Rat).SetFloat64(duration.Minutes())
	factor := new(big.Rat).SetInt64((u.Object.UsageGranularity / time.Minute).Nanoseconds())
	quantity := new(big.Rat).Quo(minutes, factor)
	ceil := new(big.Rat).SetInt(ratCeil(quantity))
	return u.SetQuantity(ceil)
}

// SetStartEnd sets the start date and the end date of an Usage
func (u *Usage) SetStartEnd(start, end time.Time) error {
	roundedStart := start.Round(u.Object.UsageGranularity)
	if roundedStart.After(start) {
		roundedStart = roundedStart.Add(-u.Object.UsageGranularity)
	}
	roundedEnd := end.Round(u.Object.UsageGranularity)
	if roundedEnd.Before(end) {
		roundedEnd = roundedEnd.Add(u.Object.UsageGranularity)
	}
	return u.SetDuration(roundedEnd.Sub(roundedStart))
}

// BillableQuantity returns the billable quantity of an Usage
func (u *Usage) BillableQuantity() *big.Rat {
	if u.Quantity.Cmp(ratZero) < 1 {
		return big.NewRat(0, 1)
	}

	//return math.Ceil(u.Quantity/u.Object.UnitQuantity) * u.Object.UnitQuantity
	quantityQuotient := new(big.Rat).Quo(u.Quantity, u.Object.UnitQuantity)
	ceil := new(big.Rat).SetInt(ratCeil(quantityQuotient))
	return new(big.Rat).Mul(ceil, u.Object.UnitQuantity)
}

// LostQuantity returns the lost quantity of an Usage
func (u *Usage) LostQuantity() *big.Rat {
	//return u.BillableQuantity() - math.Max(u.Quantity, 0)

	return new(big.Rat).Sub(u.BillableQuantity(), ratMax(u.Quantity, ratZero))
}

// Total returns the total of an Usage
func (u *Usage) Total() *big.Rat {
	//return math.Min(u.Object.UnitPrice * u.BillableQuantity(), u.Object.UnitPriceCap)

	total := new(big.Rat).Mul(u.BillableQuantity(), u.Object.UnitPrice)
	total = total.Quo(total, u.Object.UnitQuantity)
	return ratMin(total, u.Object.UnitPriceCap)
}

// TotalString returns a string representing the total price of an Usage + its currency
func (u *Usage) TotalString() string {
	return PriceString(u.Total(), u.Object.Currency)
}
