package passwordgenerator_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
)

func TestResolve(t *testing.T) {
	got, err := passwordgenerator.Resolve(false, "keep-me")
	if err != nil || got != "keep-me" {
		t.Fatalf("Resolve(false): got %q err %v", got, err)
	}

	got, err = passwordgenerator.Resolve(true, "keep-me")
	if err != nil || got != "keep-me" {
		t.Fatalf("Resolve(true, set): got %q err %v", got, err)
	}

	got, err = passwordgenerator.Resolve(true, "")
	if err != nil || !isPassword(got) {
		t.Fatalf("Resolve(true, empty): got %q err %v", got, err)
	}
}

func TestResolvePtr(t *testing.T) {
	existing := "keep-me"
	got, err := passwordgenerator.ResolvePtr(true, &existing)
	if err != nil || got == nil || *got != "keep-me" {
		t.Fatalf("ResolvePtr existing: got %v err %v", got, err)
	}

	got, err = passwordgenerator.ResolvePtr(false, nil)
	if err != nil || got != nil {
		t.Fatalf("ResolvePtr false nil: got %v err %v", got, err)
	}

	got, err = passwordgenerator.ResolvePtr(true, nil)
	if err != nil || got == nil || !isPassword(*got) {
		t.Fatalf("ResolvePtr generate: got %v err %v", got, err)
	}
}
