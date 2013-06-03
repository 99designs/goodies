// Package config reads a JSON formatted configuration file into an arbitrary struct
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func Load(path string, c interface{}) {
	Parse(Read(path), c)
}

func Read(path string) []byte {
	configJson, openerr := ioutil.ReadFile(path)
	if openerr != nil {
		log.Fatal("Could not open config file: ", openerr)
	}
	return configJson
}

func Parse(configJson []byte, c interface{}) {
	parseError := json.Unmarshal(configJson, c)

	if parseError != nil {
		log.Fatal("Could not parse config file: ", parseError)
	}
}
