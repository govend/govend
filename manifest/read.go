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

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	ext := filepath.Ext(file)

	switch ext {
	case ".json":
		if err := json.Unmarshal(bytes, v); err != nil {
			return err
		}
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("vendor manifest file type %s not supported", ext)
	}

	// save format for when we write the file back to disk
	format = ext[1:]

	return nil
}
