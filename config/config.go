// Package config provides helper functions for reading a
// JSON formatted configuration file into an arbitrary struct
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/mtibben/gocase"
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

type MarshalledEnvironmentVar struct {
	EnvKey          string
	EnvVal          string
	StructFieldName string
	Error           error
}

// LoadFromEnv iterates through a
// struct's fields and tries to find matching
// environment variables.
// Returns a map of environment key and values that were
// successfully set into the struct
func LoadFromEnv(v interface{}, prefix string) (result []MarshalledEnvironmentVar) {
	pointerValue := reflect.ValueOf(v)
	structValue := pointerValue.Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)

		if fieldValue.CanSet() {
			envKey := strings.ToUpper(prefix) + gocase.ToUpperSnake(structField.Name)
			envVal := os.Getenv(envKey)

			if envVal != "" {
				// create a json blob with the env data
				jsonStr := ""
				if fieldValue.Kind() == reflect.String {
					jsonStr = fmt.Sprintf(`{"%s": "%s"}`, structField.Name, envVal)
				} else {
					jsonStr = fmt.Sprintf(`{"%s": %s}`, structField.Name, envVal)
				}

				err := json.Unmarshal([]byte(jsonStr), v)
				result = append(result, MarshalledEnvironmentVar{envKey, envVal, structField.Name, err})
			}
		}
	}

	return
}
