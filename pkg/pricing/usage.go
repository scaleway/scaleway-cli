package pricing

import "math"

type Usage struct {
	PricingObject *PricingObject
	Quantity      float64
}

func NewUsageByPath(objectPath string, quantity float64) Usage {
	return NewUsage(CurrentPricing.GetByPath(objectPath), quantity)
}

func NewUsage(object *PricingObject, quantity float64) Usage {
	return Usage{
		PricingObject: object,
		Quantity:      quantity,
	}
}

func (u *Usage) BillableQuantity() float64 {
	return math.Ceil(math.Max(u.Quantity, 0)/u.PricingObject.UnitQuantity) * u.PricingObject.UnitQuantity
}

func (u *Usage) LostQuantity() float64 {
	return u.BillableQuantity() - math.Max(u.Quantity, 0)
}

func (u *Usage) Total() float64 {
	total := u.PricingObject.UnitPrice * u.BillableQuantity()
	return math.Min(total, u.PricingObject.UnitPriceCap)
}
