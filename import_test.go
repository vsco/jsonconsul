package jsonconsul

func ExampleJsonImport_Run() {
	ji := &JsonImport{Filename: "example.json"}
	ji.Run()

	je := &JsonExport{Prefix: "foo", IncludePrefix: true, JsonValues: true}
	je.Run()

	// Output:
	// {"foo":{"bar":"test","blah":"Test","bool":true,"do":"TEST","float":1.23,"loud":{"asd":{"bah":"test"}},"null":null}}
}
