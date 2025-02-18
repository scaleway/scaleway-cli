package passwordgenerator_test

import (
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
)

const (
	minLength = 8
)

func TestPassword(t *testing.T) {
	password, err := passwordgenerator.GeneratePassword(21, 1, 1, 1, 1)
	if err != nil {
		t.Fatalf("Error generating password: %v", err)
	}
	if !isPassword(password) {
		t.Fatalf("Generated password does not respect constraints")
	}

	// Test different lengths
	for i := minLength; i < 30; i++ {
		password, err = passwordgenerator.GeneratePassword(i, 1, 1, 1, 1)
		if err != nil {
			t.Fatalf("Error generating password: %v", err)
		}
		if !isPassword(password) {
			t.Fatalf("Generated password does not respect constraints for length %d", i)
		}
	}

	// Test different minimum character types
	password, err = passwordgenerator.GeneratePassword(21, 5, 5, 5, 5)
	if err != nil {
		t.Fatalf("Error generating password: %v", err)
	}
	if !isPassword(password) {
		t.Fatalf("Generated password does not respect constraints")
	}
}

func isPassword(s string) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSymbol bool
	hasLength := len(s)

	for _, value := range s {
		switch {
		case value >= '0' && value <= '9':
			hasNumber = true
		case value >= 'A' && value <= 'Z':
			hasUpperCase = true
		case value >= 'a' && value <= 'z':
			hasLowercase = true
		case strings.Contains(passwordgenerator.SpecialSymbols, string(value)): //nolint:gocritic
			hasSymbol = true
		}
	}

	return hasNumber && hasUpperCase && hasLowercase && hasSymbol && hasLength >= minLength
}
