package extension

import (
	"testing"
)

func TestCreatePackagedExtension(t *testing.T) {
	filepath, err := CreatePackagedExtension(`{ 
		"name" : "custom.jmx.ipa.jvm",
		"version": "1.0.0"
	}
	`)

	println(filepath)

	if(err != nil) {
		t.Fatalf("fail %v", err)
	}
}