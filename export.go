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
	Poll bool
	// How frequently should we poll the consul server for
	// changes. This should be in seconds
	PollFrequency time.Duration

	FlattenedKVs map[string]interface{}

	currentJson []byte
}

func (c *JsonExport) ParseFlags(args []string) {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&c.Prefix, "prefix", "", "What KV prefix should I track?")
	flags.StringVar(&c.ConfigFile, "config", "", "Place to output the config file. Default is config.json")
	flags.BoolVar(&c.Timestamp, "timestamp", false, "Should I create timestamped values of this")
	flags.BoolVar(&c.Poll, "poll", false, "Should I poll for changes")

	frequency := flags.Int("poll-frequency", 60, "How frequently should we poll the consul agent. In seconds")
	c.PollFrequency = time.Duration(*frequency)

	flags.Parse(args)
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

func (c *JsonExport) GenerateJson() []byte {
	c.FlattenedKVs = consulToNestedMap(c.Prefix)

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
		fmt.Println("Waiting", time.Second*c.PollFrequency)
		select {
		case <-time.After(time.Second * c.PollFrequency):
			c.Run()
		}
	}
}
