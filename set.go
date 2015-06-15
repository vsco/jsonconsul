package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type JsonSet struct {
	// The Consul Key to set.
	Key string
	// The value that was passed by the command line.
	value string
	// The actual Json Value that is going to be saved to Consul.
	JsonValue []byte
}

func (js *JsonSet) ParseFlags(args []string) {
	if len(args) < 2 {
		fmt.Println("jsonconsul set full/key/path \"true\"")
		os.Exit(-1)
	}

	js.Key = args[0]
	js.value = args[1]
}

func (js *JsonSet) lintedJson() []byte {
	var (
		unmarshalled interface{}
	)

	err := json.Unmarshal([]byte(js.value), &unmarshalled)
	if err != nil {
		log.Fatal("Can't set the key ", js.Key, " invalid value: ", js.value)
	}

	return []byte(js.value)
}

func (js *JsonSet) Run() {
	js.JsonValue = js.lintedJson()

	setConsulKV(js.Key, js.JsonValue)
}
