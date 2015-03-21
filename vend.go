package main

//go:generate go run mkpkgs.go

import (
	"fmt"
	"os"
)

// vend vendors packages into the vendor directory.
func vend(verbose bool) error {

	// default value for vendors yaml file.
	vendFile = "vendor/vend.yml"

	// check if vend file exists.
	if _, err := os.Stat(file); err == nil {

		// check verbosity
		if verbose {
			fmt.Println("vend.yml found: processing...")
		}

		// read the vendors file.
		err, vendors := readVendFile(vendFile)
		if err != nil {
			return err
		}
	}

	// iterate through each .go file

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
}
