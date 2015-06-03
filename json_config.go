package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type JsonConfig struct {
	// KV Path in Consul
	Prefix string
	// Place to put the config file
	ConfigFile string
	// If we are overwriting the ConfigFile should we timestamp
	// the versions so that there is a trail.
	Timestamp bool
	// Should we poll for consul changes.
	Poll bool
	// How frequently should we poll the consul server for
	// changes. This should be in seconds
	PollFrequency time.Duration

	currentJson []byte
}

func (c *JsonConfig) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&c.Prefix, "prefix", "", "What KV prefix should I track?")
	flags.StringVar(&c.ConfigFile, "config", "", "Place to output the config file. Default is config.json")
	flags.BoolVar(&c.Timestamp, "timestamp", false, "Should I create timestamped values of this")
	flags.BoolVar(&c.Poll, "poll", false, "Should I poll for changes")

	frequency := flags.Int("poll_frequency", 60, "How frequently should we poll the consul agent. In seconds")
	c.PollFrequency = time.Duration(*frequency)

	flags.Parse(os.Args[1:])

}

func (c *JsonConfig) getConsulToMap() (v map[string]interface{}) {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV() // Lookup the pair

	pairs, _, err := kv.List(c.Prefix, nil)
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

func (c *JsonConfig) fileNameWithTimestamp() string {
	return fmt.Sprintf("%s.%d", c.ConfigFile, int32(time.Now().Unix()))
}

func (c *JsonConfig) WriteFile(newJson []byte) {
	if bytes.Equal(c.currentJson, newJson) {
		// File didn't change.
		return
	}

	fileName := c.ConfigFile
	if c.Timestamp {
		fileName = c.fileNameWithTimestamp()
	}

	err := ioutil.WriteFile(fileName, newJson, 0644)
	if err != nil {
		panic(err)
	}

	// Symlink files
	if c.Timestamp {
		os.Symlink(fileName, c.ConfigFile)
	}

	c.currentJson = newJson
}

func (c *JsonConfig) GenerateJson() []byte {
	consulMap := c.getConsulToMap()

	js, err := json.Marshal(consulMap)
	if err != nil {
		panic(err)
	}

	return js
}

func (c *JsonConfig) Run() {
	json := c.GenerateJson()

	if c.ConfigFile == "" {
		fmt.Println(string(json))
	} else {
		c.WriteFile(json)
	}
}

func (c *JsonConfig) RunWatcher() {
	for {
		fmt.Println("Waiting", time.Second*c.PollFrequency)
		select {
		case <-time.After(time.Second * c.PollFrequency):
			c.Run()
		}
	}
}
