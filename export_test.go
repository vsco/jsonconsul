package jsonconsul

func ExampleJsonExport_Run_IncludePrefix() {
	ji := &JsonImport{Filename: "example.json"}
	ji.Run()

	config := &JsonExport{
		Prefix:        "foo",
		IncludePrefix: true,
		JsonValues:    true,
	}
	config.Run()

	// Output:
	// {"foo":{"bar":"test","blah":"Test","bool":true,"do":"TEST","float":1.23,"loud":{"asd":{"bah":"test"}},"null":null}}
}

func ExampleJsonExport_Run_NoIncludePrefix() {
	ji := &JsonImport{Filename: "example.json"}
	ji.Run()

	config := &JsonExport{
		Prefix:        "foo",
		IncludePrefix: false,
		JsonValues:    true,
	}
	config.Run()

	// Output:
	// {"bar":"test","blah":"Test","bool":true,"do":"TEST","float":1.23,"loud":{"asd":{"bah":"test"}},"null":null}
}

func ExampleJsonExport_Run_IncludePrefixNoJsonValues() {
	ji := &JsonImport{Filename: "example.json"}
	ji.Run()

	config := &JsonExport{
		Prefix:        "foo",
		IncludePrefix: true,
		JsonValues:    false,
	}
	config.Run()

	// Output:
	// {"foo":{"bar":"\"test\"","blah":"\"Test\"","bool":"true","do":"\"TEST\"","float":"1.23","loud":{"asd":{"bah":"\"test\""}},"null":"null"}}
}
