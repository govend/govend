package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const defaultVendorFile = "vendor/vend.yml"

// vend vendors packages into the vendor directory.
func vend(verbose bool) error {

	// verbosity
	if verbose {
		fmt.Print("determining project path...")
	}

	// determine the project import path
	projectpath, err := importpath(".")
	if err != nil {
		return err
	}

	// verbosity
	if verbose {
		fmt.Println(" 			" + projectpath)
	}

	// verbosity
	if verbose {
		fmt.Print("scanning for external packages...")
	}

	// scan for external packages
	pkgs, err := scan(".")
	if err != nil {
		return err
	}

	// filter out packages internal to the project
	pkgs = rmprefix(projectpath, pkgs)

	// verbosity
	if verbose {
		fmt.Println(" 		" + strconv.Itoa(len(pkgs)) + " packages found")
	}

	// set default value for vendors yaml file.
	vendpath := defaultVendorFile

	// verbosity
	if verbose {
		if vendpath == defaultVendorFile {
			fmt.Println("setting default vend file...			" + vendpath)
		} else {
			fmt.Println("setting custom vend file... 			" + vendpath)
		}
	}

	// create an empty slice of vendors to fill.
	vf := new(vendors)

	// verbosity
	if verbose {
		fmt.Print("checking the vend file exist...")
	}

	// check if vend file path exists.
	if _, err := os.Stat(vendpath); err == nil {

		// verbosity
		if verbose {
			fmt.Println("			file exists")
		}

		// verbosity
		if verbose {
			fmt.Println("processing " + vendpath + "...")
		}

		// read the vendors file.
		if err := load(vendpath, vf); err != nil {
			return err
		}

	} else {

		// verbosity
		if verbose {
			fmt.Println("			file missing")
		}
	}

	log.Fatal(pkgs)

	// iterate through each .go file

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
	return nil
}
