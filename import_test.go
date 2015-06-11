package main

import (
	"github.com/hashicorp/consul/api"
)

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
