package passwordgenerator

import (
	crypto "crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const (
	numbers        = "0123456789"
	lowerLetters   = "abcdedfghijklmnopqrstuvwxyz"
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialSymbols = "!$%^&*()_+{}:@[];'#<>?,./|\\\\-=?"
	allSet         = lowerLetters + upperLetters + specialSymbols + numbers
)

func GeneratePassword(length, minNumbers, minLower, minUpper, minSymbol int) (string, error) {
	if length < (minNumbers + minLower + minUpper + minSymbol) {
		return "", errors.New("length is less than the sum of minNumbers, minLower, minUpper, and minSymbol")
	}

	var password strings.Builder

	for i := 0; i < minNumbers; i++ {
		random, err := randInt(len(numbers))
		if err != nil {
			return "", err
		}
		password.WriteString(string(numbers[random]))
	}

	for i := 0; i < minUpper; i++ {
		random, err := randInt(len(upperLetters))
		if err != nil {
			return "", err
		}
		password.WriteString(string(upperLetters[random]))
	}

	for i := 0; i < minLower; i++ {
		random, err := randInt(len(lowerLetters))
		if err != nil {
			return "", err
		}
		password.WriteString(string(lowerLetters[random]))
	}

	for i := 0; i < minSymbol; i++ {
		random, err := randInt(len(specialSymbols))
		if err != nil {
			return "", err
		}
		password.WriteString(string(specialSymbols[random]))
	}

	remainingLength := length - minNumbers - minLower - minUpper - minSymbol
	for i := 0; i < remainingLength; i++ {
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

func randInt(len int) (int, error) {
	i, err := crypto.Int(crypto.Reader, big.NewInt(int64(len)))
	if err != nil {
		return 0, err
	}
	return int(i.Int64()), nil
}
