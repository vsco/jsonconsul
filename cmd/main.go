package main

import (
	"fmt"
	"github.com/vsco/jsonconsul"
	"log"
	"os"
)

const (
	Name  = "jsonconsul"
	usage = `
Usage: %s [mode] [options]

Mode:

  watch    Watch for changes in Consul and generate json files.
  export   Export the keys as a nested JSON file.
  import   Import json file into appropriate KV pairs in Consul.
  set      Set a K/V in Consul to linted JSON value.
`
)

func showUsage() {
	fmt.Printf(usage, Name)
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(-1)
	}

	switch os.Args[1] {
	case "watch":
		jsonExport := &jsonconsul.JsonExport{Watch: true}
		jsonExport.ParseFlags(os.Args[2:])
		jsonExport.RunWatcher()
	case "export":
		jsonExport := &jsonconsul.JsonExport{Watch: false}
		jsonExport.ParseFlags(os.Args[2:])
		err := jsonExport.Run()
		if err != nil {
			log.Fatalln(err)
		}
	case "import":
		jsonImport := &jsonconsul.JsonImport{}
		jsonImport.ParseFlags(os.Args[2:])
		err := jsonImport.Run()
		if err != nil {
			log.Fatalln(err)
		}
	case "set":
		jsonSet := &jsonconsul.JsonSet{}
		jsonSet.ParseFlags(os.Args[2:])
		err := jsonSet.Run()
		if err != nil {
			log.Fatalln(err)
		}
	default:
		showUsage()
	}
}
