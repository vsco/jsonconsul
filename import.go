package main

import (
	"encoding/json"
	//	"flag"
	"fmt"
	"io/ioutil"
	//	"os"
	"github.com/hashicorp/consul/api"
)

type JsonImport struct {
	// Prefix to load the config under. If empty then loads to the
	// root kv node.
	Prefix string
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

func (ji *JsonImport) prefixedKey(key string) string {
	return fmt.Sprintf("%s/%s", ji.Prefix, key)
}

func (ji *JsonImport) setConsulValues() {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV() // Lookup the pair

	for k, v := range ji.flattened {
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

	ji.flattened = make(map[string]string)
	ji.keysFromJson(unmarshalled, "")
	ji.setConsulValues()
}
