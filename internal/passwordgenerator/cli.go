package passwordgenerator

import "github.com/scaleway/scaleway-cli/v2/core"

const defaultLength = 21

// Generate returns a 21-character password (mixed case, digits, symbols).
func Generate() (string, error) {
	return GeneratePassword(defaultLength, 1, 1, 1, 1)
}

// ArgSpec returns the shared generate-password CLI flag.
func ArgSpec(defaultEnabled bool) *core.ArgSpec {
	defaultValue := "false"
	if defaultEnabled {
		defaultValue = "true"
	}

	return &core.ArgSpec{
		Name:       "generate-password",
		Short:      `Will generate a 21 character-length password that contains a mix of upper/lower case letters, numbers and special symbols`,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter(defaultValue),
	}
}

// Resolve returns password, generating one when generate is set and password is empty.
func Resolve(generate bool, password string) (string, error) {
	if !generate || password != "" {
		return password, nil
	}

	return Generate()
}

// ResolvePtr returns password, generating one when generate is set and password is nil/empty.
func ResolvePtr(generate bool, password *string) (*string, error) {
	if password != nil && *password != "" {
		return password, nil
	}
	if !generate {
		return password, nil
	}

	pwd, err := Generate()
	if err != nil {
		return nil, err
	}

	return &pwd, nil
}
