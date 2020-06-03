package core

// Example represents an example for the usage of a CLI command.
type Example struct {

	// Short is the title given to the example.
	Short string

	// ArgsJSON is a JSON encoded representation of the request used in the example. Only one of ArgsJSON or Raw should be provided.
	ArgsJSON string

	// Raw is a raw example. Only one of ArgsJSON or Raw should be provided.
	Raw string
}
