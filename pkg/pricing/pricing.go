package pricing

import (
	"math/big"
	"time"
)

// Object represents a Pricing item definition
type Object struct {
	Path             string
	Identifier       string
	Currency         string
	UsageUnit        string
	UnitPrice        *big.Rat
	UnitQuantity     *big.Rat
	UnitPriceCap     *big.Rat
	UsageGranularity time.Duration
}

// List represents a list of Object
type List []Object

// CurrentPricing tries to be up-to-date with the real pricing
// we cannot guarantee of these values since we hardcode values for now
// later, we should be able to call a dedicated pricing API
var CurrentPricing List

func init() {
	CurrentPricing = List{
		{
			Path:             "/compute/c1/run",
			Identifier:       "aaaaaaaa-aaaa-4aaa-8aaa-111111111112",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(2, 1000),     // 0.002
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(1000, 1000),  // 1
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/compute/c2s/run",
			Identifier:       "aaaaaaaa-aaaa-4aaa-8aaa-222222222222",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(20, 1000),    // 0.02
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(10000, 1000), // 10
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/compute/c2m/run",
			Identifier:       "aaaaaaaa-aaaa-4aaa-8aaa-444444444444",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(32, 1000),    // 0.032
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(16000, 1000), // 16
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/compute/c2l/run",
			Identifier:       "aaaaaaaa-aaaa-4aaa-8aaa-333333333333",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(44, 1000),    // 0.044
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(22000, 1000), // 22
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/compute/vc1s/run",
			Identifier:       "cccccccc-6ab1-4131-a35e-000000000001",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(2, 1000),     // 0.002
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(1000, 1000),  // 1
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/compute/vc1m/run",
			Identifier:       "cccccccc-6ab1-4131-a35e-000000000002",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(6, 1000),     // 0.006
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(3000, 1000),  // 3
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/compute/vc1l/run",
			Identifier:       "cccccccc-6ab1-4131-a35e-000000000003",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(10, 1000),    // 0.01
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(5000, 1000),  // 5
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/ip/dynamic",
			Identifier:       "467116bf-4631-49fb-905b-e07701c21111",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(2, 1000),     // 0.002
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(990, 1000),   // 0.99
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/ip/reserved",
			Identifier:       "467116bf-4631-49fb-905b-e07701c22222",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(2, 1000),     // 0.002
			UnitQuantity:     big.NewRat(60000, 1000), // 60
			UnitPriceCap:     big.NewRat(990, 1000),   // 0.99
			UsageGranularity: time.Minute,
		},
		{
			Path:             "/storage/local/ssd/storage",
			Identifier:       "bbbbbbbb-bbbb-4bbb-8bbb-111111111144",
			Currency:         "EUR",
			UnitPrice:        big.NewRat(2, 1000),     // 0.002
			UnitQuantity:     big.NewRat(50000, 1000), // 50
			UnitPriceCap:     big.NewRat(1000, 1000),  // 1
			UsageGranularity: time.Hour,
		},
	}
}

// GetByPath returns an object matching a path
func (pl *List) GetByPath(path string) *Object {
	for _, object := range *pl {
		if object.Path == path {
			return &object
		}
	}
	return nil
}

// GetByIdentifier returns an object matching a identifier
func (pl *List) GetByIdentifier(identifier string) *Object {
	for _, object := range *pl {
		if object.Identifier == identifier {
			return &object
		}
	}
	return nil
}
