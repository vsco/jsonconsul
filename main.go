package main

import (
	"flag"
	"time"
)

var (
	config *JsonConfig
)

func init() {
	config = &JsonConfig{}

	flag.StringVar(&config.Prefix, "prefix", "", "What KV prefix should I track?")
	flag.StringVar(&config.ConfigFile, "config", "", "Place to output the config file. Default is config.json")
	flag.BoolVar(&config.Timestamp, "timestamp", false, "Should I create timestamped values of this")
	flag.BoolVar(&config.Poll, "poll", false, "Should I poll for changes")

	frequency := flag.Int("poll_frequency", 60, "How frequently should we poll the consul agent. In seconds")
	config.PollFrequency = time.Duration(*frequency)
}

func main() {
	flag.Parse()

	if config.Poll {
		config.Poller()
	} else {
		config.Handler()
	}
}
