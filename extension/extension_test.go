package extension

import (
	"archive/zip"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePackagedExtension(t *testing.T) {
	packagedExtension, err := CreatePackagedExtension(`{ 
		"name" : "custom.jmx.ipa.jvm",
		"version": "1.0.0"
	}
	`)
	if(err != nil) {
		t.Fatalf("fail %v", err)
	}

	defer packagedExtension.Dispose()

	reader, err := zip.OpenReader(packagedExtension.FilePath)
	if(err != nil) {
		t.Fatalf("fail %v", err)
	}
	defer reader.Close()
	
	f, err := reader.Open("plugin.json")
	if(err != nil) {
		t.Fatalf("Fail %v", err)
	}

	bytes, err := io.ReadAll(f)
	if(err != nil) {
		t.Fatalf("Fail %v", err)
	}

	var result map[string]string
	err = json.Unmarshal(bytes, &result)
	if(err != nil) {
		t.Fatalf("Fail %v", err)
	}

	assert.Equal(t, "custom.jmx.ipa.jvm", result["name"], "name mismatch")
	assert.Equal(t, "1.0.0", result["version"], "version mismatch")


}

