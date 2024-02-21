package extension

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ExtensionTestDataFile = "../testdata/custom.jmx.testext.json"

type DynatraceExtensionTest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Version string `json:"version"`
}

func TestCreatePackagedExtension(t *testing.T) {
	file, _ := os.Open(ExtensionTestDataFile)
	payload, _ := io.ReadAll(file)

	packagedExtension, err := CreatePackagedExtension(string(payload))
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

	var result DynatraceExtensionTest
	err = json.Unmarshal(bytes, &result)
	if(err != nil) {
		t.Fatalf("Fail %v", err)
	}

	assert.Equal(t, "custom.jmx.testext", result.Name, "name mismatch")
	assert.Equal(t, "1.0", result.Version, "version mismatch")


}

