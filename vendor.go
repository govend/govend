package main

// vendor describes a vendored dependecy.
// vendor contains the import path and revision hash.
type vendor struct {
	Path string `json:"path" yaml:"path"`
	Rev  string `json:"rev,"omitempty"" yaml:"rev,"omitempty""`
}
