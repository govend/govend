package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var format string

// Read reads a vend.yml file and returns an array of vendors.
func Read(file string, v *[]Vendor) error {

	// read the vend file into bytes
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// get the file extension
	ext := filepath.Ext(file)

	// unmarshal by file extension
	switch ext {

	// unmarshal JSON
	case ".json":
		if err := json.Unmarshal(bytes, v); err != nil {
			return err
		}

	// unmarshal YML
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

	// error on unsupported extension type
	default:
		return fmt.Errorf("vendor manifest file type %s not supported", ext)
	}

	// save format for when we write the file back to disk
	format = ext[1:]

	return nil
}
