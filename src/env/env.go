package env

import (
	"os"
	"reflect"
	"strconv"
)

// Get reads the environment and returns a populated, verified instance of Variables.
func Get() (Variables, error) {
	v := Variables{}

	for envVar, fieldName := range evMapping {
		value := os.Getenv(envVar)

		field := reflect.ValueOf(&v).Elem().FieldByName(fieldName)
		switch field.Type().Kind().String() {
		case "string":
			field.SetString(value)
		case "bool":
			field.SetBool(parseBool(value))
		case "int":
			if value == "" {
				field.SetInt(0)
				continue
			}

			t, e := strconv.Atoi(value)
			if e != nil {
				panic(e)
			}

			field.SetInt(int64(t))
		}
	}

	if err := v.Validate(); err != nil {
		return Variables{}, err
	}

	return v, nil
}

func parseBool(str string) bool {
	b, e := strconv.ParseBool(str)
	if e != nil {
		return false
	}

	return b
}
