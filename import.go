package main

import (
	"encoding/json"
	//	"flag"
	"fmt"
	"io/ioutil"
	//	"os"
)

type JsonImport struct {
	// File containing the Json
	Filename string

	flattened map[string]string
}

func (ji *JsonImport) ParseFlags(args []string) {

}

func (ji *JsonImport) keysFromJson(nested interface{}, prefix string) {
	for k, v := range nested.(map[string]interface{}) {
		newPrefix := fmt.Sprintf("%s/%s", prefix, k)
		switch v.(type) {
		case string:
			ji.flattened[newPrefix] = v.(string)
		case interface{}:
			ji.keysFromJson(v, newPrefix)
		}
	}
}

func (ji *JsonImport) Run() {
	var unmarshalled map[string]interface{}

	ji.flattened = make(map[string]string)

	fileOutput, err := ioutil.ReadFile(ji.Filename)
	if err != nil {
		panic(err)
	}

	// Import the json file
	if err := json.Unmarshal(fileOutput, &unmarshalled); err != nil {
		panic(err)
	}

	ji.keysFromJson(unmarshalled, "")

	// Set the consul values
}
