package main

func ExampleJsonExport_Run_IncludePrefix() {
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
	config := &JsonExport{
		Prefix:        "foo",
		IncludePrefix: true,
		JsonValues:    false,
	}
	config.Run()

	// Output:
	// {"foo":{"bar":"\"test\"","blah":"\"Test\"","bool":"true","do":"\"TEST\"","float":"1.23","loud":{"asd":{"bah":"\"test\""}},"null":"null"}}
}
