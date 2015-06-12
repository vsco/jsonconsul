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

	je := &JsonExport{Prefix: "foo", IncludePrefix: true, JsonValues: true}
	je.Run()

	// Output:
	// {"foo":{"bar":"test","blah":"Test","bool":true,"do":"TEST","float":1.23,"loud":{"asd":{"bah":"test"}},"null":null}}
}
