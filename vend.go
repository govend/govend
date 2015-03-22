package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const defaultVendorFile = "vendor/vend.yml"

// vend vendors packages into the vendor directory.
func vend(verbose bool) error {

	// determine the absolute file path for the provided directory
	currentpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// set default value for vendors yaml file.
	vendpath := defaultVendorFile

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
		fmt.Print("scanning for external unvendored packages...")
	}

	// scan for external packages
	pkgs, err := scan(".")
	if err != nil {
		return err
	}

	// filter out unvendored packages
	uvpkgs := rmprefix(projectpath, pkgs)

	// verbosity
	if verbose {
		fmt.Println(" 		" + strconv.Itoa(len(uvpkgs)) + " packages found")
	}

	// filter out vendored packages
	vpkgs := pickprefix(projectpath+"/vendor/", pkgs)

	// check if no externally vendored or unvendored packages exist
	if len(uvpkgs) < 1 && len(vpkgs) < 1 {

		// get stats on the path
		if _, err := os.Stat(filepath.Base(vendpath)); err != nil {

			// check if the path does not exist
			if os.IsNotExist(err) {
				return nil
			}

			return err
		}

		// remove everthing in the vendor directory
		os.RemoveAll(filepath.Base(vendpath))

		return nil
	}

	// check vpkgs is not empty
	if len(vpkgs) > 0 {

		// iterate over vpkgs
		for _, pkg := range vpkgs {

			// remove project path to create a complete absolute filepath
			vpath := pkg[len(projectpath):]

			// get stats on the pkg
			if _, err := os.Stat(filepath.Join(currentpath, vpath)); err != nil {

				// check if the path does not exist
				if os.IsNotExist(err) {

					// verbosity
					if verbose {
						fmt.Println("missing vendored code for " + pkg)
					}

					// clean pkg path to be unvendored
					pkg = pkg[len(projectpath+"/vendor/"):]

					// append package into the unvendored package object
					uvpkgs = append(uvpkgs, pkg)
				}

				return err
			}
		}
	}

	// check uvpkgs is not empty
	if len(uvpkgs) > 0 {

		// iterate over uvpkgs
		for _, pkg := range uvpkgs {

			fmt.Println(pkg)
		}

	}

	log.Fatalln("done")

	// create an empty slice of vendors to fill.
	var vf vendors

	// check if vend file path exists.
	if _, err := os.Stat(vendpath); err == nil {

		// verbosity
		if verbose {
			fmt.Println("loading " + vendpath + "...")
		}

		// read the vendors file.
		if err := load(vendpath, &vf); err != nil {
			return err
		}

		// check if the vend file is empty
		if len(vf) < 1 {

			// remove everthing in the vendor directory
			os.RemoveAll(filepath.Base(vendpath))
		}

	} else {

		// verbosity
		if verbose {
			fmt.Println("			file missing")
		}
	}

	log.Fatal(pkgs)

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
	return nil
}
