package session

import (
	"testing"
)

func TestRoundTrip(t *testing.T) {
	sd := sessionData{
		Id:    generateSessionId(),
		Store: make(map[string]string),
	}
	sd.Store["TEST"] = "Fabulous"

	roundTrip, err := unmarshalSessionData(marshalSessionData(sd))
	if err != nil {
		t.Fatal(err)
	}
	if roundTrip.Store["TEST"] != "Fabulous" {
		t.Fatal("Failed to roundtrip")
	}

}
