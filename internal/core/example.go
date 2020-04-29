package core

// Example represents an example for the usage of a CLI command.
type Example struct {

	// Short is the title given to the example.
	Short string

	// ArgJSON is a JSON encoded representation of the arguments used in the example.
	// Only one of ArgJSON or Raw should be provided.
	ArgJSON string

	// Raw is a raw example.
	// Only one of ArgJSON or Raw should be provided.
	Raw string
}
