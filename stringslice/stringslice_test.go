package stringslice

import (
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	result := Contains("12345678", []string{"abc", "123"})
	expected := true
	if result != expected {
		t.Errorf("Test 1: expected %b, got %b", expected, result)
	}

	result = Contains("hgfd", []string{"abc", "234"})
	expected = false
	if result != expected {
		t.Errorf("Test 2: expected %b, got %b", expected, result)
	}
}

func TestRemoveStringsFromString(t *testing.T) {
	result := RemoveStringsFromString("1abc2345678a123bc", []string{"abc", "123"})
	expected := "45678"
	if result != expected {
		t.Errorf("Test 1: expected %s, got %s", expected, result)
	}
}

func TestMap(t *testing.T) {
	result := Map([]string{"ABC", "Def", "ghi"}, strings.ToLower)
	expected := []string{"abc", "def", "ghi"}

	if !Equal(result, expected) {
		t.Errorf("Test 1: expected %s, got %s", expected, result)
	}
}

func TestEqual(t *testing.T) {
	result := Equal([]string{"ABC", "Def", "ghi"}, []string{"ABC", "Def", "ghi"})
	expected := true

	if result != expected {
		t.Errorf("Test 1: expected %b, got %b", expected, result)
	}

	result = Equal([]string{"ABC", "Def", "ghi"}, []string{"ABC", "Deg", "ghi"})
	expected = false

	if result != expected {
		t.Errorf("Test 2: expected %b, got %b", expected, result)
	}

	result = Equal([]string{"ABC", "Def", "ghi"}, []string{"ABC", "Def", "ghi", "asdf"})
	expected = false

	if result != expected {
		t.Errorf("Test 2: expected %b, got %b", expected, result)
	}
}
