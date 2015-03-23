package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/go/vcs"
)

const defaultVendorFile = "_vendor/vend.yml"

type person struct {
	Name string `json:"name" yaml:"name"` // Supporting both JSON and YAML.
	Age  int    `json:"age" yaml:"age"`
}

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
		fmt.Println(" 	" + strconv.Itoa(len(uvpkgs)) + " packages found")
	}

	// filter out vendored packages
	vpkgs := pickprefix(projectpath+"/_vendor/", pkgs)

	// check if no externally vendored or unvendored packages exist
	if len(uvpkgs) < 1 && len(vpkgs) < 1 {

		// get stats on the path
		if _, err := os.Stat(filepath.Base(vendpath)); err != nil {

			// check if the error is not that the path does not exist
			if !os.IsNotExist(err) {
				return err
			}
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
					pkg = pkg[len(projectpath+"/_vendor/"):]

					// append package into the unvendored package object
					uvpkgs = append(uvpkgs, pkg)
				}

				return err
			}
		}
	}

	// create an empty slice of vendors to fill.
	var vf []vendor

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

			if verbose {
				fmt.Println("			empty file")
			}

			// remove the vend file
			os.Remove(vendpath)
		}

	} else {

		// verbosity
		if verbose {
			fmt.Println("			file missing")
			fmt.Println(vendpath)
		}
	}

	// check uvpkgs is not empty
	if len(uvpkgs) > 0 {

		// create a repo map of package paths to RepoRoots
		rmap := make(map[string]*vcs.RepoRoot)

		// iterate over uvpkgs
		// remove package imports that might already be included
		// example: "gopkg.in/mgo.v2/bson" -> "gopkg.in/mgo.v2"
		for _, pkg := range uvpkgs {

			// determine import path dynamically by pinging repository
			r, err := vcs.RepoRootForImportDynamic(pkg, false)
			if err != nil {
				return err
			}

			// check if package path is missing from repo map
			if _, ok := rmap[pkg]; !ok {

				// add the RepoRoot to the repo map
				rmap[pkg] = r
			}
		}

		// check that the repo map is not empty
		if len(rmap) > 0 {

			// iterate through the rmap
			for _, r := range rmap {

				// create a directory for the pkg
				os.MkdirAll(filepath.Dir("_vendortemp/"+r.Root), 0777)

				for _, v := range vf {

					if r.Root == v.Path && len(v.Rev) > 0 {
						// create the pkg
						r.VCS.CreateAtRev("_vendortemp/"+r.Root, r.Repo, v.Rev)
						fmt.Println("Root: " + r.Root + " | Repo " + r.Repo + " | " + v.Rev)
						goto VendorMatch
					}
				}
				fmt.Println("Root: " + r.Root + " | Repo " + r.Repo + " | no rev")
				r.VCS.Create("_vendortemp/"+r.Root, r.Repo)

			VendorMatch:
			}
		}

		log.Fatal(rmap)
	}

	// os.RemoveAll("_vendor")
	// CopyDir("_vendortemp", "_vendor")
	// os.RemoveAll("_vendortemp")

	fmt.Println("")
	log.Fatalln(pkgs)

	//    check if package is in vendors file

	//      check by tag versions if we should update

	//    if not in vendor file then add it to vendors
	return nil
}
