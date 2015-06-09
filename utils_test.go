package main

import "testing"

func TestWordify(t *testing.T) {
	actual := wordify("Hello World 42 !!")
	expected := "Hello_World_42"
	if actual != expected {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}
}
