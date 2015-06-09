package main

import "testing"

func TestWordify(t *testing.T) {
	actual := wordify("Hello World 42 !!")
	expected := "Hello_World_42"
	if actual != expected {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}
}

func TestTruncIf(t *testing.T) {
	actual := truncIf("Hello World", 5, false)
	expected := "Hello World"
	if actual != expected {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}

	actual = truncIf("Hello World", 5, true)
	expected = "Hello"
	if actual != expected {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}

	actual = truncIf("Hello World", 50, false)
	expected = "Hello World"
	if actual != expected {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}

	actual = truncIf("Hello World", 50, true)
	expected = "Hello World"
	if actual != expected {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}
}

func TestPathToSCPPathparts(t *testing.T) {
	dir, base := PathToSCPPathparts("/etc/passwd")
	expected := []string{"/etc", "passwd"}
	actual := []string{dir, base}
	if actual[0] != expected[0] || actual[1] != expected[1] {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}

	dir, base = PathToSCPPathparts("/etc")
	expected = []string{"/", "etc"}
	actual = []string{dir, base}
	if actual[0] != expected[0] || actual[1] != expected[1] {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}

	dir, base = PathToSCPPathparts("/etc/")
	expected = []string{"/", "etc"}
	actual = []string{dir, base}
	if actual[0] != expected[0] || actual[1] != expected[1] {
		t.Errorf("returned value is invalid [actual: %s][expected: %s]", actual, expected)
	}
}
