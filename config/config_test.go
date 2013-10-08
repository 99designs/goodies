package config

import (
	"testing"
)

type ApplicationConfig struct {
	BindAddress              string
	DatabaseConnectionString string
	MoreStuff                []NestedConfig
}

type NestedConfig struct {
	Foo string
	Bar string
}

func TestParse(t *testing.T) {
	sampleConfig := `
{
  "BindAddress": ":8008",
  "DatabaseConnectionString": "dbname/user/pass",
  "MoreStuff": [
    {
      "Foo": "bar",
      "Bar": "baz"
    }
  ]
}
  `

	var c ApplicationConfig
	Parse([]byte(sampleConfig), &c)

	if c.BindAddress != ":8008" {
		t.Error("Unexpected value in BindAddress")
	}
	if c.DatabaseConnectionString != "dbname/user/pass" {
		t.Error("Unexpected value in DatabaseConnectionString")
	}
	if len(c.MoreStuff) != 1 {
		t.Error("Unexpected length of MoreStuff")
	}
	if c.MoreStuff[0].Foo != "bar" {
		t.Error("Unexpected value in c.MoreStuff[0].Foo")
	}
	if c.MoreStuff[0].Bar != "baz" {
		t.Error("Unexpected value in c.MoreStuff[0].Bar")
	}
}
