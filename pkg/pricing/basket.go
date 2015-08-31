package pricing

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
