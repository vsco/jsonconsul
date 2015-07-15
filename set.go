package jsonconsul

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const (
	Undefined int = iota
	Bool
	Number
	String
	Array
	Object
)

type JsonSet struct {
	// The Consul Key to set.
	Key string
	// The value that was passed by the command line.
	Value string
	// What we expect the value type to be.
	ExpectedType int
	// The actual Json Value that is going to be saved to Consul.
	JsonValue []byte
	// Is the value going to be json.
	IsJsonValue bool
}

func (js *JsonSet) setExpectedType(t string) {
	switch t {
	case "bool":
		js.ExpectedType = Bool
	case "number":
		js.ExpectedType = Number
	case "int":
		js.ExpectedType = Number
	case "float":
		js.ExpectedType = Number
	case "string":
		js.ExpectedType = String
	case "array":
		js.ExpectedType = Array
	case "object":
		js.ExpectedType = Object
	default:
		js.ExpectedType = Undefined
	}
}

func (js *JsonSet) expectedType() string {
	switch js.ExpectedType {
	case Bool:
		return "bool"
	case Number:
		return "number"
	case String:
		return "string"
	case Array:
		return "array"
	case Object:
		return "object"
	}

	return "undefined"
}

func (js *JsonSet) checkExpectedType(unmarshalled interface{}) error {
	if js.ExpectedType == Undefined {
		return nil
	}
	switch unmarshalled.(type) {
	case bool:
		if js.ExpectedType != Bool {
			return fmt.Errorf("Invalid type. Value is a bool. Expected %s", js.expectedType())
		}
	case float64:
		if js.ExpectedType != Number {
			return fmt.Errorf("Invalid type. Value is a number. Expected %s", js.expectedType())
		}
	case string:
		if js.ExpectedType != String {
			return fmt.Errorf("Invalid type. Value is a string. Expected %s", js.expectedType())
		}
	case []interface{}:
		if js.ExpectedType != Array {
			return fmt.Errorf("Invalid type. Value is an array. Expected %s", js.expectedType())
		}
	case map[string]interface{}:
		if js.ExpectedType != Object {
			return fmt.Errorf("Invalid type. Value is an object. Expected %s", js.expectedType())
		}
	}

	return nil
}

func (js *JsonSet) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	expectedType := flags.String("expected-type", "undefined", "What is the expected type for the value? bool, int, float, string, object, array")
	flags.BoolVar(&js.IsJsonValue, "json-value", true, "Is the value going to be set as json")

	flags.Parse(args)

	js.setExpectedType(*expectedType)

	leftovers := flags.Args()
	if len(leftovers) < 2 {
		fmt.Println("jsonconsul set full/key/path \"true\"")
		os.Exit(-1)
	}

	js.Key = leftovers[0]
	js.Value = leftovers[1]
}

func (js *JsonSet) lintedJson() ([]byte, error) {
	var (
		unmarshalled interface{}
	)

	err := json.Unmarshal([]byte(js.Value), &unmarshalled)
	if err != nil {
		return nil, fmt.Errorf("Can't set the key %s invalid value: %s", js.Key, js.Value)
	}

	err = js.checkExpectedType(unmarshalled)
	if err != nil {
		return nil, err
	}

	return []byte(js.Value), nil
}

func (js *JsonSet) Run() error {
	var (
		err error
		val []byte
	)

	if js.IsJsonValue {
		js.JsonValue, err = js.lintedJson()
		if err != nil {
			return err
		}

		val = js.JsonValue
	} else {
		val = []byte(js.Value)
	}

	err = setConsulKV(js.Key, val)
	if err != nil {
		return err
	}

	return nil
}
