package config

import (
	e "github.com/99designs/goodies/test/expect"
	"testing"
)

type ApplicationConfig struct {
	BindAddress              string
	DatabaseConnectionString string
	Clients                  []ConfigClient
}

type ConfigClient struct {
	Name string
	Key  string
}

func TestParse(t *testing.T) {
	sampleConfig := `
{
  "BindAddress": ":8008",
  "DatabaseConnectionString": "contests/auth/auth",
  "Clients": [
    {
      "Name": "contests",
      "Key": "ContestsAuthDevelopmentPSK"
    }
  ]
}
  `
	var c ApplicationConfig
	Parse([]byte(sampleConfig), &c)
	e.Expect(t, "Bind address", c.BindAddress, ":8008")
	e.Expect(t, "DB connection string", c.DatabaseConnectionString, "contests/auth/auth")
	e.Expect(t, "Number of clients", len(c.Clients), 1)
	e.Expect(t, "Client URI", c.Clients[0].Key, "ContestsAuthDevelopmentPSK")
	e.Expect(t, "Client URI", c.Clients[0].Name, "contests")
}
