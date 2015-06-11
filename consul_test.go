package main

import (
	"testing"
)

func Test_consulPrefixedKey(t *testing.T) {
	if consulPrefixedKey("foobar", "blah") != "foobar/blah" {
		t.Error("Did not return correct key")
	}

	if consulPrefixedKey("", "/blah") == "foobar/blah" {
		t.Error("Did not return correct key")
	}
}

func TestinterfaceToConsulFlattenMap(t *testing.T) {

}

func TestconsulFlattenMap(t *testing.T) {

}

func TestconsulToNestedMap(t *testing.T) {

}
