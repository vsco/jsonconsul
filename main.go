package main

import (
	"fmt"
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
		jsonExport := &JsonExport{}
		jsonExport.RunWatcher()
	case "export":
		jsonExport := &JsonExport{}
		jsonExport.Run()
	case "import":
		jsonImport := &JsonImport{}
		jsonImport.ParseFlags(os.Args[1:])
		jsonImport.Run()
	default:
		showUsage()
	}
}
