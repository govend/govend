package main

//go:generate go run mkpkgs.go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// vend vendors packages into the vendor directory.
func vend(verbose bool) error {

	// default value for vendors yaml file.
	vendpath = "vendor/vend.yml"

	// check if vend file path exists.
	if _, err := os.Stat(vendpath); err == nil {

		// check verbosity
		if verbose {
			fmt.Println("vend.yml found: processing...")
		}

		// read the vendors file.
		vendors, err := loadVendors(vendpath)
		if err != nil {
			return err
		}
	}

	// iterate through each .go file

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
}

// loadVendors reads a vend.yml file and returns an array of vendors.
func loadVendors(filepath string) (vendors, error) {

	// read the vend file into bytes.
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// create an empty slice of vendors to fill.
	v := new(vendors)

	// get the file extension.
	ext := filepath.Ext(filename)

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

	return v, nil
}
