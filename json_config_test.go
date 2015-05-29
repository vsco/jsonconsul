package main

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func SetupKeys() {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV() // Lookup the pair

	for _, i := range []string{"foo/bar", "foo/blah", "foo/do", "foo/loud/asd/bah"} {
		p := &api.KVPair{Key: i, Value: []byte("test")}
		_, err := kv.Put(p, nil)
		if err != nil {
			panic(err)
		}
	}

}

func TestJsonConfig_Generate(t *testing.T) {
	SetupKeys()

	config := &JsonConfig{
		Prefix: "foo",
	}
	config.GenerateJson()

	// Output:
	// {"foo":{"bar":"test","blah":"test","do":"test","loud":{"asd":{"bah":"test"}}}}
}
