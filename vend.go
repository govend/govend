package main

import (
	"fmt"
	"os"
)

func vend(verbose bool) error {

	file = "vendor/vend.yml"

	// create vendors to work with.
	vendors := NewVendors()

	// check if vend file exists.
	if _, err := os.Stat(file); err == nil {
		fmt.Println("vend.yml found: processing...")
		if err := vendors.Read("./vendor/vend.yml"); err != nil {
			return err
		}
	}

	// iterate through each .go file

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
}
