// Methods which interact directly with Consul. Should isolate the
// code for consul here.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
)

var (
	client *api.Client
	kv     *api.KV
)

func interfaceToConsulFlattenedMap(nested interface{}, prefix string, output map[string]interface{}) {
	for k, v := range nested.(map[string]interface{}) {
		newPrefix := fmt.Sprintf("%s/%s", prefix, k)

		switch v.(type) {
		case map[string]interface{}:
			interfaceToConsulFlattenedMap(v, newPrefix, output)
		default:
			output[newPrefix] = v
		}
	}
}

func consulToFlattenedMap(prefix string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	nested, err := consulToNestedMap(prefix, true)
	if err != nil {
		return nil, err
	}

	interfaceToConsulFlattenedMap(nested, prefix, output)

	return output, nil
}

func consulToNestedMap(prefix string, includePrefix bool) (v map[string]interface{}, err error) {
	v = make(map[string]interface{})

	kv := client.KV() // Lookup the pair

	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return v, err
	}

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

	if !includePrefix {
		nodes := strings.Split(prefix, "/")
		for _, node := range nodes {
			switch v[node].(type) {
			case map[string]interface{}:
				v, _ = v[node].(map[string]interface{})
			}
		}
	}

	return v, nil
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

func setConsulKVs(prefix string, kvMap map[string]interface{}) error {
	for k, v := range kvMap {
		value, err := json.Marshal(v)
		if err != nil {
			return err
		}

		if err := setConsulKV(consulPrefixedKey(prefix, k), value); err != nil {
			return err
		}
	}

	return nil
}

func setConsulKV(key string, value []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: value,
	}

	_, err := kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	var (
		err error
	)

	client, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalln(err)
	}
	kv = client.KV()
}
