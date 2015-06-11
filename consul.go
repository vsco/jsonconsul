// Methods which interact directly with Consul. Should isolate the
// code for consul here.
package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strings"
)

func interfaceToConsulFlattenedMap(nested interface{}, prefix string, output map[string]string) {
	for k, v := range nested.(map[string]interface{}) {
		newPrefix := fmt.Sprintf("%s/%s", prefix, k)
		switch v.(type) {
		case string:
			output[newPrefix] = v.(string)
		case interface{}:
			interfaceToConsulFlattenedMap(v, newPrefix, output)
		}
	}
}

func consulToFlattenedMap(prefix string) map[string]string {
	output := make(map[string]string)
	nested := consulToNestedMap(prefix)
	interfaceToConsulFlattenedMap(nested, prefix, output)

	return output
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

func consulPrefixedKey(prefix, key string) (newKey string) {
	if prefix == "" {
		newKey = key
	} else {
		newKey = fmt.Sprintf("%s/%s", prefix, strings.TrimPrefix(key, "/"))
	}

	newKey = strings.TrimPrefix(newKey, "/")

	return
}

func setConsulKVs(prefix string, kvMap map[string]string) {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV()

	for k, v := range kvMap {
		p := &api.KVPair{
			Key:   consulPrefixedKey(prefix, k),
			Value: []byte(v),
		}

		_, err := kv.Put(p, nil)
		if err != nil {
			panic(err)
		}
	}
}
