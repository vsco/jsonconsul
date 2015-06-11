package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type JsonImport struct {
	// Prefix to load the config under. If empty then loads to the
	// root kv node.
	Prefix string
	// File containing the Json to be converted to KVs.
	Filename string

	FlattenedKVs map[string]string
}

func (ji *JsonImport) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&ji.Prefix, "prefix", "", "What prefix should the Key Values be stored under.")
	flags.StringVar(&ji.Filename, "json-file", "", "Json file that will be imported into Consul.")
	flags.Parse(args)

	if ji.Filename == "" {
		fmt.Println("Include filename with -json-file")
		os.Exit(-1)
	}
}

func (ji *JsonImport) readFile() (unmarshalled map[string]interface{}) {
	fileOutput, err := ioutil.ReadFile(ji.Filename)
	if err != nil {
		panic(err)
	}

	// Import the json file
	if err := json.Unmarshal(fileOutput, &unmarshalled); err != nil {
		panic(err)
	}

	return
}

func (ji *JsonImport) Run() {
	ji.FlattenedKVs = make(map[string]string)
	unmarshalled := ji.readFile()

	interfaceToConsulFlattenedMap(unmarshalled, "", ji.FlattenedKVs)
	setConsulKVs(ji.Prefix, ji.FlattenedKVs)
}
