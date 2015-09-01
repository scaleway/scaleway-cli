package pricing

import (
	"math/big"
	"time"
)

type Usage struct {
	PricingObject *PricingObject
	Quantity      *big.Rat
}

func NewUsageByPath(objectPath string) Usage {
	return NewUsageByPathWithQuantity(objectPath, 0)
}

func NewUsageByPathWithQuantity(objectPath string, quantity float64) Usage {
	return NewUsageWithQuantity(CurrentPricing.GetByPath(objectPath), quantity)
}

func NewUsageWithQuantity(object *PricingObject, quantity float64) Usage {
	return Usage{
		PricingObject: object,
		Quantity:      new(big.Rat).SetFloat64(quantity),
	}
}

func NewUsage(object *PricingObject) Usage {
	return NewUsageWithQuantity(object, 0)
}

func (u *Usage) SetQuantity(quantity *big.Rat) error {
	u.Quantity = ratMax(quantity, ratZero)
	return nil
}

func (u *Usage) SetDuration(duration time.Duration) error {
	minutes := new(big.Rat).SetFloat64(duration.Minutes())
	factor := new(big.Rat).SetInt64((u.PricingObject.UsageGranularity / time.Minute).Nanoseconds())
	quantity := new(big.Rat).Quo(minutes, factor)
	ceil := new(big.Rat).SetInt(ratCeil(quantity))
	return u.SetQuantity(ceil)
}

func (u *Usage) SetStartEnd(start, end time.Time) error {
	roundedStart := start.Round(u.PricingObject.UsageGranularity)
	if roundedStart.After(start) {
		roundedStart = roundedStart.Add(-u.PricingObject.UsageGranularity)
	}
	roundedEnd := end.Round(u.PricingObject.UsageGranularity)
	if roundedEnd.Before(end) {
		roundedEnd = roundedEnd.Add(u.PricingObject.UsageGranularity)
	}
	return u.SetDuration(roundedEnd.Sub(roundedStart))
}

func (u *Usage) BillableQuantity() *big.Rat {
	if u.Quantity.Cmp(big.NewRat(0, 1)) < 1 {
		return big.NewRat(0, 1)
	}

	//return math.Ceil(u.Quantity/u.PricingObject.UnitQuantity) * u.PricingObject.UnitQuantity
	unitQuantity := new(big.Rat).SetFloat64(u.PricingObject.UnitQuantity)
	quantityQuotient := new(big.Rat).Quo(u.Quantity, unitQuantity)
	ceil := new(big.Rat).SetInt(ratCeil(quantityQuotient))
	return new(big.Rat).Mul(ceil, unitQuantity)
}

func (u *Usage) LostQuantity() *big.Rat {
	//return u.BillableQuantity() - math.Max(u.Quantity, 0)

	return new(big.Rat).Sub(u.BillableQuantity(), ratMax(u.Quantity, ratZero))
}

func (u *Usage) Total() *big.Rat {
	//return math.Min(u.PricingObject.UnitPrice * u.BillableQuantity(), u.PricingObject.UnitPriceCap)

	unitPrice := new(big.Rat).SetFloat64(u.PricingObject.UnitPrice)
	total := new(big.Rat).Mul(u.BillableQuantity(), unitPrice)

	unitPriceCap := new(big.Rat).SetFloat64(u.PricingObject.UnitPriceCap)
	return ratMin(total, unitPriceCap)
}
