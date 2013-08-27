// Package stringslice implements simple functions to manipulate string slices
package stringslice

import (
	"strings"
)

// Contains returns true if any string in substrs is within s
func Contains(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}

// RemoveStringsFromString removes any string in substrs from s
func RemoveStringsFromString(s string, substrs []string) string {
	for loop := true; loop; {
		loop = false
		for _, substr := range substrs {
			lastS := s
			s = strings.Join(strings.Split(s, substr), "")
			if lastS != s {
				loop = true
				break
			}
		}
	}

	return s
}

// MapStringSlice returns a slice containing all the elements of ss
// after applying the callback function to each one
func Map(ss []string, callback func(string) string) []string {
	newStrings := make([]string, len(ss))
	for i, s := range ss {
		newStrings[i] = callback(s)
	}

	return newStrings
}

// Equal returns a boolean reporting whether a == b
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
