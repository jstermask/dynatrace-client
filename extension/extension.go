package extension

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/jstermask/go-dynatrace-client/model"
)

const FolderPattern string = "dynatrace_extension"
const InnerFileName string = "plugin.json"


func CreatePackagedExtension(payload string) (*string, error) {
	var metadata model.DynatraceExtensionMetadata
	err := json.Unmarshal([]byte(payload), &metadata)
	if err != nil {
		return nil, err
	}

	// create a zip file containing a plugin.json file with payload content
	zipDir, err := os.MkdirTemp(os.TempDir(), FolderPattern)
	if err != nil {
		return nil, err
	}

	zipFilePath := fmt.Sprintf("%s/%s.zip", zipDir, metadata.Name)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	entry, err := zipWriter.Create(InnerFileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(entry, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	return &zipFilePath, nil

}
