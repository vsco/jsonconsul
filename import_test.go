package main

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestJsonImport_prefixedKey(t *testing.T) {
	ji := &JsonImport{Prefix: "foobar"}
	if ji.prefixedKey("blah") != "foobar/blah" {
		t.Error("Did not return correct key")
	}

	ji.Prefix = "/foobar"
	if ji.prefixedKey("blah") != "foobar/blah" {
		t.Error("Did not return correct key")
	}
	if ji.prefixedKey("/blah") != "foobar/blah" {
		t.Error("Did not return correct key")
	}
}

func ExampleJsonImport_Run() {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV() // Lookup the pair
	kv.DeleteTree("", nil)

	ji := &JsonImport{Filename: "example.json"}
	ji.Run()

	je := &JsonExport{Prefix: "foo"}
	je.Run()

	// Output:
	// {"foo":{"bar":"test","blah":"Test","do":"TEST","loud":{"asd":{"bah":"test"}}}}
}
