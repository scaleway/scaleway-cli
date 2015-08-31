package pricing

type PricingObject struct {
	Path             string
	Identifier       string
	Currency         string
	UnitPrice        float64
	UnitQuantity     float64
	UnitPriceCap     float64
	UsageUnit        string
	UsageGranularity string
}

type PricingList []PricingObject

// CurrentPricing tries to be up-to-date with the real pricing
// we cannot guarantee of these values since we hardcode values for now
// later, we should be able to call a dedicated pricing API
var CurrentPricing PricingList

func init() {
	CurrentPricing = PricingList{
		{
			Path:             "/compute/c1/run",
			Identifier:       "aaaaaaaa-aaaa-4aaa-8aaa-111111111111",
			Currency:         "EUR",
			UnitPrice:        0.012,
			UnitQuantity:     60,
			UnitPriceCap:     6,
			UsageGranularity: "minute",
		},
		{
			Path:             "/ip/dynamic",
			Identifier:       "467116bf-4631-49fb-905b-e07701c2db11",
			Currency:         "EUR",
			UnitPrice:        0.004,
			UnitQuantity:     60,
			UnitPriceCap:     1.99,
			UsageGranularity: "minute",
		},
		{
			Path:             "/ip/reserved",
			Identifier:       "467116bf-4631-49fb-905b-e07701c2db22",
			Currency:         "EUR",
			UnitPrice:        0.004,
			UnitQuantity:     60,
			UnitPriceCap:     1.99,
			UsageGranularity: "minute",
		},
		{
			Path:             "/storage/local/ssd/storage",
			Identifier:       "bbbbbbbb-bbbb-4bbb-8bbb-111111111113",
			Currency:         "EUR",
			UnitPrice:        0.004,
			UnitQuantity:     50,
			UnitPriceCap:     2,
			UsageGranularity: "hour",
		},
	}
}

// GetByPath returns an object matching a path
func (pl *PricingList) GetByPath(path string) *PricingObject {
	for _, object := range *pl {
		if object.Path == path {
			return &object
		}
	}
	return nil
}

// GetByIdentifier returns an object matching a identifier
func (pl *PricingList) GetByIdentifier(identifier string) *PricingObject {
	for _, object := range *pl {
		if object.Identifier == identifier {
			return &object
		}
	}
	return nil
}
