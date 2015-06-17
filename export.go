package jsonconsul

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

func (c *JsonExport) WriteFile(newJson []byte) error {
	if bytes.Equal(c.currentJson, newJson) {
		// File didn't change.
		return nil
	}

	fileName := c.ConfigFile
	if c.Timestamp {
		fileName = c.fileNameWithTimestamp()
	}

	err := ioutil.WriteFile(fileName, newJson, 0644)
	if err != nil {
		return err
	}

	// Symlink files
	if c.Timestamp {
		err = os.Symlink(fileName, c.ConfigFile)
		if err != nil {
			return err
		}
	}

	c.currentJson = newJson

	return nil
}

func (c *JsonExport) jsonifyValues(kvs map[string]interface{}) error {
	for k, v := range kvs {
		switch v.(type) {
		case string:
			var jsonVal interface{}
			err := json.Unmarshal([]byte(v.(string)), &jsonVal)
			if err != nil {
				return err
			}
			kvs[k] = jsonVal
		case map[string]interface{}:
			if err := c.jsonifyValues(v.(map[string]interface{})); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *JsonExport) GenerateJson() ([]byte, error) {
	var (
		err error
	)

	c.FlattenedKVs, err = consulToNestedMap(c.Prefix, c.IncludePrefix)
	if err != nil {
		return nil, err
	}

	if c.JsonValues {
		err = c.jsonifyValues(c.FlattenedKVs)
		if err != nil {
			return nil, err
		}
	}

	js, err := json.Marshal(c.FlattenedKVs)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func (c *JsonExport) Run() error {
	json, err := c.GenerateJson()
	if err != nil {
		return err
	}

	if c.ConfigFile == "" {
		fmt.Println(string(json))
	} else {
		err := c.WriteFile(json)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *JsonExport) RunWatcher() {
	for {
		err := c.Run()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("Waiting", time.Second*c.WatchFrequency)
		<-time.After(time.Second * c.WatchFrequency)
	}
}
