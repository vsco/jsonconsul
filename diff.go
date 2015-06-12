package main

import (
	"flag"
	"fmt"
	// "github.com/sergi/go-diff/diffmatchpatch"
	"os"
)

type JsonDiff struct {
	// Prefix to load the config under. If empty then loads to the
	// root kv node.
	Prefix string
	// File containing the Json to be converted to KVs.
	Filename string
	// Are the KVs json values?
	JsonValues bool
}

func (ji *JsonDiff) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&ji.Prefix, "prefix", "", "What prefix should the Key Values be stored under.")
	flags.BoolVar(&ji.JsonValues, "json-values", true, "Have the values that are returned by Consul be parsed as json.")
	flags.Parse(args)

	leftovers := flags.Args()
	if len(leftovers) == 0 {
		fmt.Println("Include json file to compare")
		os.Exit(-1)
	} else {
		ji.Filename = leftovers[0]
	}
}

func (ji *JsonDiff) Run() {
	// Not yet implemented
	// config := &JsonExport{
	// 	Prefix:        "foo",
	// 	IncludePrefix: true,
	// 	JsonValues:    false,
	// }
	// config.Run()

	// jsonImport = JsonImport{}
	// unmarshalled := ji.readFile()

	// Get the JsonImport values to convert to key values
	// ji.flattened = make(map[string]string)
	// ji.keysFromJson(unmarshalled, "")

	// Get the JsonExport values to get the keys values
	// Get the export values

	// Do a deep check of the values.

}
