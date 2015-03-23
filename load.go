package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// load reads a vend.yml file and returns an array of vendors.
func load(file string, v *[]vendor) error {

	// read the vend file into bytes.
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// get the file extension.
	ext := filepath.Ext(file)

	// unmarshal by file extension.
	switch ext {

	// unmarshal JSON.
	case ".json":
		if err := json.Unmarshal(bytes, v); err != nil {
			return err
		}

	// unmarshal YML.
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

	// error on unsupported extension type.
	default:
		return errors.New("vend file type " + ext + "not supported")
	}

	return nil
}
