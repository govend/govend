package main

// vendor describes a vendored dependecy.
// vendor contains the import path and revision hash.
type vendor struct {
	path string `yaml:"path"`
	rev  string `yaml:"rev"`
}

// vendors is a slice of vendor.
type vendors []vendor
