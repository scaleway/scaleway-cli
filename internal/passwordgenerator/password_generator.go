package passwordgenerator

import (
	crypto "crypto/rand"
	"errors"
	"math/big"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
)

const (
	numbers        = "0123456789"
	lowerLetters   = "abcdedfghijklmnopqrstuvwxyz"
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SpecialSymbols = "!$%^&*()_+{}:@[];'#<>?,./|\\\\-=?"
	allSet         = lowerLetters + upperLetters + SpecialSymbols + numbers
)

func GeneratePassword(length, minNumbers, minLower, minUpper, minSymbol int) (string, error) {
	if length < (minNumbers + minLower + minUpper + minSymbol) {
		return "", errors.New(
			"length is less than the sum of minNumbers, minLower, minUpper, and minSymbol",
		)
	}

	var password strings.Builder

	for range minNumbers {
		random, err := randInt(len(numbers))
		if err != nil {
			return "", err
		}
		password.WriteString(string(numbers[random]))
	}

	for range minUpper {
		random, err := randInt(len(upperLetters))
		if err != nil {
			return "", err
		}
		password.WriteString(string(upperLetters[random]))
	}

	for range minLower {
		random, err := randInt(len(lowerLetters))
		if err != nil {
			return "", err
		}
		password.WriteString(string(lowerLetters[random]))
	}

	for range minSymbol {
		random, err := randInt(len(SpecialSymbols))
		if err != nil {
			return "", err
		}
		password.WriteString(string(SpecialSymbols[random]))
	}

	remainingLength := length - minNumbers - minLower - minUpper - minSymbol
	for range remainingLength {
		random, err := randInt(len(allSet))
		if err != nil {
			return "", err
		}
		password.WriteString(string(allSet[random]))
	}
	inRune := []rune(password.String())
	for i := range inRune {
		j, err := randInt(i + 1)
		if err != nil {
			return "", err
		}
		inRune[i], inRune[j] = inRune[j], inRune[i]
	}

	return string(inRune), nil
}

func randInt(length int) (int, error) {
	i, err := crypto.Int(crypto.Reader, big.NewInt(int64(length)))
	if err != nil {
		return 0, err
	}

	return int(i.Int64()), nil
}

// GenerateIfNeeded stores a generated 21-character password in *password when
// generate is true and *password is empty, returning the generated password
// (or "" when none was generated).
func GenerateIfNeeded(generate bool, password *string) (string, error) {
	if !generate || *password != "" {
		return "", nil
	}

	generated, err := GeneratePassword(21, 1, 1, 1, 1)
	if err != nil {
		return "", err
	}
	*password = generated

	return generated, nil
}

// GeneratePasswordArgSpec is the shared --generate-password flag spec.
func GeneratePasswordArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:       "generate-password",
		Short:      `Will generate a 21 character-length password that contains a mix of upper/lower case letters, numbers and special symbols`,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("true"),
	}
}
