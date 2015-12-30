package pricing

import (
	"math/big"
	"time"
)

// Basket represents a billing basket
type Basket []Usage

// NewBasket return a new instance of Basket
func NewBasket() *Basket {
	return &Basket{}
}

// Add adds an Usage to a Basket
func (b *Basket) Add(usage Usage) error {
	*b = append(*b, usage)
	return nil
}

// Length returns the amount of Usages in a Basket
func (b *Basket) Length() int {
	return len(*b)
}

// SetDuration sets the duration for each Usages of a Basket
func (b *Basket) SetDuration(duration time.Duration) error {
	var err error
	for i, usage := range *b {
		err = usage.SetDuration(duration)
		if err != nil {
			return err
		}
		(*b)[i] = usage
	}
	return nil
}

// Total computes the Usage.Total() of all the Usages of a Basket
func (b *Basket) Total() *big.Rat {
	total := new(big.Rat)
	for _, usage := range *b {
		total = total.Add(total, usage.Total())
	}
	return total
}
