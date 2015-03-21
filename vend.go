package main

import (
	"fmt"
	"os"
)

// vend vendors packages into the vendor directory.
func vend(verbose bool) error {

	// default value for vendors yaml file.
	vendpath := "vendor/vend.yml"

	// create an empty slice of vendors to fill.
	v := new(vendors)

	// check if vend file path exists.
	if _, err := os.Stat(vendpath); err == nil {

		// check verbosity
		if verbose {
			fmt.Println("vend.yml found: processing...")
		}

		// read the vendors file.
		if err := load(vendpath, v); err != nil {
			return err
		}

	} else {
		fmt.Println(vendpath + " not found")
		pkgs, err := scan(".")
		if err != nil {
			return err
		}
		fmt.Print(pkgs)
	}

	// iterate through each .go file

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
	return nil
}
