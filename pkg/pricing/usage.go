package pricing

import "math/big"

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

func (u *Usage) SetQuantity(quantity *big.Rat) {
	u.Quantity = quantity
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
