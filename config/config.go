// Package config provides helper functions for reading a
// JSON formatted configuration file into an arbitrary struct
package config

import (
	"encoding/json"
	"io/ioutil"
)

// Load reads the JSON-encoded file and marshalls
// the contents into the value pointed at by v.
// Panics if unsuccessful
func Load(filename string, v interface{}) {
	Parse(read(filename), v)
}

func read(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	return data
}

// Parse parses the JSON-encoded data and stores the result
// into the value pointed at by c.
// Panics if unsuccessful
func Parse(jsondata []byte, v interface{}) {
	err := json.Unmarshal(jsondata, v)
	if err != nil {
		panic(err.Error())
	}
}
