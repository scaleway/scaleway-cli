package pricing

import "math/big"

type Basket []Usage

func NewBasket() *Basket {
	return &Basket{}
}

func (b *Basket) Add(usage Usage) error {
	*b = append(*b, usage)
	return nil
}

func (b *Basket) Length() int {
	return len(*b)
}

func (b *Basket) Total() *big.Rat {
	total := new(big.Rat)
	for _, usage := range *b {
		total = total.Add(total, usage.Total())
	}
	return total
}
