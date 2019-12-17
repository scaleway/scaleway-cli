package core

// Example represents an example for the usage of a CLI command.
type Example struct {

	// Short is the title given to the example.
	Short string

	// Request is a JSON encoded representation of the request used in the example. Only one of Request or Raw should be provided.
	Request string

	// Raw is a raw example. Only one of Request or Raw should be provided.
	Raw string
}
