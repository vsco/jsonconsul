package main

import (
	"fmt"
	"testing"
)

func ExampleJsonSet_RunBadJson() {
	ji := &JsonSet{}
	ji.ParseFlags([]string{"blah/blah", "\"a"})
	err := ji.Run()
	fmt.Println(err)

	// Output:
	// Can't set the key blah/blah invalid value: "a
}

func ExampleJsonSet_RunBadExpectedType() {
	ji := &JsonSet{}
	ji.ParseFlags([]string{"blah/blah", "true"})
	ji.setExpectedType("int")
	err := ji.Run()
	fmt.Println(err)

	// Output:
	// Invalid type. Value is a bool. Expected number
}

func TestJsonSet_RunGoodExpectedType(t *testing.T) {
	ji := &JsonSet{}
	ji.ParseFlags([]string{"blah/blah", "true"})
	ji.setExpectedType("bool")
	err := ji.Run()

	if err != nil {
		t.Error(err)
	}
}
