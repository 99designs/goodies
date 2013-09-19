package session

import (
	"fmt"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	sd := sessionData{
		Id:    generateSessionId(),
		Store: make(map[string]string),
	}

	sd.Set("TEST", "Fabulous")
	result := sd.Get("TEST")
	if result != "Fabulous" {
		t.Errorf("Expected Fabulous, got %s", result)
	}

	data, err := sd.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	newSd := sessionData{}
	err = newSd.Marshal(data)
	if err != nil {
		t.Error(err)
	}

	result = newSd.Get("TEST")
	if result != "Fabulous" {
		t.Errorf("Expected Fabulous, got %s", result)
	}

	newSd.Clear()

	result = newSd.Get("TEST")
	if result != "Fabulous" {
		t.Errorf("Expected Fabulous, got %s", result)
	}
}
