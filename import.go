package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"os"
	"strings"
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

func (ji *JsonImport) keysFromJson(nested interface{}, prefix string) {
	for k, v := range nested.(map[string]interface{}) {
		newPrefix := fmt.Sprintf("%s/%s", prefix, k)
		switch v.(type) {
		case string:
			ji.FlattenedKVs[newPrefix] = v.(string)
		case interface{}:
			ji.keysFromJson(v, newPrefix)
		}
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

func (ji *JsonImport) prefixedKey(key string) (newKey string) {
	if ji.Prefix == "" {
		newKey = key
	} else {
		newKey = fmt.Sprintf("%s/%s", ji.Prefix, strings.TrimPrefix(key, "/"))
	}

	newKey = strings.TrimPrefix(newKey, "/")

	return
}

func (ji *JsonImport) setConsulValues() {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV() // Lookup the pair

	for k, v := range ji.FlattenedKVs {
		p := &api.KVPair{
			Key:   ji.prefixedKey(k),
			Value: []byte(v),
		}

		_, err := kv.Put(p, nil)
		if err != nil {
			panic(err)
		}
	}
}

func (ji *JsonImport) Run() {
	unmarshalled := ji.readFile()

	ji.FlattenedKVs = make(map[string]string)
	ji.keysFromJson(unmarshalled, "")
	ji.setConsulValues()
}
