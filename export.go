package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type JsonExport struct {
	// KV Path in Consul
	Prefix string
	// Place to put the config file
	ConfigFile string
	// If we are overwriting the ConfigFile should we timestamp
	// the versions so that there is a trail.
	Timestamp bool
	// Should we poll for consul changes.
	Watch bool
	// How frequently should we poll the consul server for
	// changes. This should be in seconds
	WatchFrequency time.Duration
	// Should the output include the nodes in the included prefix?
	IncludePrefix bool
	// Parse the Values as Json
	JsonValues bool

	FlattenedKVs map[string]interface{}

	currentJson []byte
}

func (c *JsonExport) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&c.Prefix, "prefix", "", "What KV prefix should I track?")
	flags.BoolVar(&c.Timestamp, "timestamp", false, "Should I create timestamped values of this")
	flags.BoolVar(&c.IncludePrefix, "include-prefix", true, "Should I remove the prefix values when exporting?")
	flags.BoolVar(&c.JsonValues, "json-values", true, "Have the values that are returned by Consul be parsed as json.")

	if c.Watch {
		frequency := flags.Int("poll-frequency", 60, "How frequently should we poll the consul agent. In seconds")
		c.WatchFrequency = time.Duration(*frequency)
	}

	flags.Parse(args)

	leftovers := flags.Args()
	if len(leftovers) != 0 {
		c.ConfigFile = leftovers[0]
	}
}

func (c *JsonExport) fileNameWithTimestamp() string {
	return fmt.Sprintf("%s.%d", c.ConfigFile, int32(time.Now().Unix()))
}

func (c *JsonExport) WriteFile(newJson []byte) {
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

func (c *JsonExport) jsonifyValues(kvs map[string]interface{}) {
	for k, v := range kvs {
		switch v.(type) {
		case string:
			var jsonVal interface{}
			json.Unmarshal([]byte(v.(string)), &jsonVal)
			kvs[k] = jsonVal
		case map[string]interface{}:
			c.jsonifyValues(v.(map[string]interface{}))
		}
	}
}

func (c *JsonExport) GenerateJson() []byte {
	c.FlattenedKVs = consulToNestedMap(c.Prefix, c.IncludePrefix)
	if c.JsonValues {
		c.jsonifyValues(c.FlattenedKVs)
	}

	js, err := json.Marshal(c.FlattenedKVs)
	if err != nil {
		panic(err)
	}

	return js
}

func (c *JsonExport) Run() {
	json := c.GenerateJson()

	if c.ConfigFile == "" {
		fmt.Println(string(json))
	} else {
		c.WriteFile(json)
	}
}

func (c *JsonExport) RunWatcher() {
	for {
		fmt.Println("Waiting", time.Second*c.WatchFrequency)
		<-time.After(time.Second * c.WatchFrequency)
		c.Run()
	}
}
