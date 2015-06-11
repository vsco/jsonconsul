package main

import (
	"github.com/hashicorp/consul/api"
	"strings"
)

func consulToFlattenedMap(prefix string) {

}

func consulToNestedMap(prefix string) (v map[string]interface{}) {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV() // Lookup the pair

	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		panic(err)
	}

	v = make(map[string]interface{})

	for _, n := range pairs {
		keyIter := v
		keys := strings.Split(n.Key, "/")

		for i, key := range keys {
			if i == len(keys)-1 {
				keyIter[key] = string(n.Value)
			} else {
				if _, ok := keyIter[key]; !ok {
					keyIter[key] = make(map[string]interface{})
				}
				keyIter = keyIter[key].(map[string]interface{})
			}
		}
	}

	return
}
