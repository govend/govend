package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/go/vcs"
)

// vend vendors packages into the vendor directory.
func vendcmd(verbose bool) error {

	// determine the absolute file path for the current local directory
	localpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

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

	// remove standard packages
	pkgs = removestdpkgs(pkgs)

	// find the unvendored packages by removing packages that contain the
	// projectpath as a prefix in the import path
	uvpkgs := removeprefix(projectpath, pkgs)

	// verbosity
	if verbose {
		fmt.Println(" 	" + strconv.Itoa(len(uvpkgs)) + " packages found")
	}

	// filter out vendored packages
	vpkgs := selectprefix(filepath.Join(projectpath, vendorDir+string(filepath.Separator)), pkgs)

	// check if no externally vendored or unvendored packages exist
	if len(uvpkgs) < 1 && len(vpkgs) < 1 {

		// get stats on the path
		if _, err := os.Stat(vendorDir); err != nil {

			// check if the error is not that the path does not exist
			if !os.IsNotExist(err) {
				return err
			}
		}

		// remove everthing in the vendor directory
		os.RemoveAll(vendorDir)

		return nil
	}

	// check vpkgs is not empty
	if len(vpkgs) > 0 {

		// iterate over vpkgs
		for _, pkg := range vpkgs {

			// remove project path to create a complete absolute filepath
			vpath := pkg[len(projectpath):]

			// get stats on the pkg
			if _, err := os.Stat(filepath.Join(localpath, vpath)); err != nil {

				// check if the path does not exist
				if os.IsNotExist(err) {

					// verbosity
					if verbose {
						fmt.Println("missing vendored code for " + pkg)
					}

					// clean pkg path to be unvendored
					pkg = pkg[len(filepath.Join(projectpath, vendorDir+string(filepath.Separator))):]

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
	if _, err := os.Stat(vendorFilePath); err == nil {

		// verbosity
		if verbose {
			fmt.Println("loading " + vendorFilePath + "...")
		}

		// read the vendors file.
		if err := load(vendorFilePath, &vf); err != nil {
			return err
		}

		// check if the vend file is empty
		if len(vf) < 1 {

			if verbose {
				fmt.Println("			empty file")
			}

			// remove the vend file
			os.Remove(vendorFilePath)
		}

	} else {

		// verbosity
		if verbose {
			fmt.Println("			file missing: " + vendorFilePath)
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
				e := err.Error()
				msg := "no go-import meta tags"
				if e[len(msg):] == msg {
					return errors.New(e + "\nAre you behind a proxy?")
				}
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
				os.MkdirAll(filepath.Dir(filepath.Join(vendorTempDir, r.Root)), 0777)

				// iterate over the vendors in the vendor file
				for _, v := range vf {

					// check if we have a match, and a given revision exists
					if r.Root == v.Path && len(v.Rev) > 0 {

						// create the repository at that specific revision
						r.VCS.CreateAtRev(filepath.Join(vendorTempDir, r.Root), r.Repo, v.Rev)
						goto RevMatch
					}
				}
				r.VCS.Create(filepath.Join(vendorTempDir, r.Root), r.Repo)

			RevMatch:
			}

			// iterate through the rmap
			for _, r := range rmap {
				os.RemoveAll(filepath.Join(vendorDir, r.Root))
				os.MkdirAll(filepath.Dir(filepath.Join(vendorDir, r.Root)), 0777)
				CopyDir(filepath.Join(vendorTempDir, r.Root), filepath.Join(vendorDir, r.Root))
			}

			os.RemoveAll(vendorTempDir)
		}
	}

	// if not in vendor file then add it to vendors
	return nil
}
