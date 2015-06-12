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

	FlattenedKVs map[string]interface{}
}

func (ji *JsonImport) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&ji.Prefix, "prefix", "", "What prefix should the Key Values be stored under.")
	flags.Parse(args)

	leftovers := flags.Args()
	if len(leftovers) == 0 {
		fmt.Println("Must pass a file to import")
		os.Exit(-1)
	} else {
		ji.Filename = leftovers[0]
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
	ji.FlattenedKVs = make(map[string]interface{})
	unmarshalled := ji.readFile()

	interfaceToConsulFlattenedMap(unmarshalled, "", ji.FlattenedKVs)
	setConsulKVs(ji.Prefix, ji.FlattenedKVs)
}
