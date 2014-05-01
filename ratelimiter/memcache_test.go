package ratelimiter

import "testing"

func TestNormaliseKey(t *testing.T) {

	result := normaliseKey(" foo bar	 baz 私は - <> 寿司が\nいい")
	expected := "+foo+bar%09+baz+%E7%A7%81%E3%81%AF+-+%3C%3E+%E5%AF%BF%E5%8F%B8%E3%81%8C%0A%E3%81%84%E3%81%84"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestNormaliseKeyLength(t *testing.T) {
	akey := ""
	for i := 0; i <= 252; i++ {
		akey += "X"
	}

	result := len(normaliseKey(akey))
	expected := 250

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
