package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Vend vendors packages into the vendor directory.
// It achieves this by following the following steps:
//
// Step 1. Identify all relative file paths necessary for the current project.
// Step 2. Identify all types of packages currently present in the project.
// Step 3. If the vendors.yml manifest file exists, load it in memory.
// Step 4. Verify vendored packages and treat bad ones as unvendored packages.
// Step 5. Identify package repositories and filter out repo subpackages.
// Step 6. Download and vendor packages.
// Step 7. Write the vendors.yml manifest file.
// Step 8. Rewrite import paths.
//
func vendcmd(verbose bool) error {

	//
	// Step 1. Identify all relative file paths necessary for the current project.
	//
	////

	// determine the absolute file path for the current local directory
	localpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// determine the project import path
	projectImportPath, err := importpath(".")
	if err != nil {
		return err
	}

	// define vendor directory paths
	vendorProjectPath := filepath.Join(projectImportPath, vendorPath)
	vendorProjectPathSlashed := vendorProjectPath + string(filepath.Separator)

	// verbosity
	if verbose {
		fmt.Print("identifying project paths... 			complete")
		fmt.Print("scanning for external unvendored packages...")
	}

	//
	// Step 2. Identify all types of packages currently present in the project.
	//
	////

	// scan for external packages
	pkgs, err := scan(".")
	if err != nil {
		return err
	}

	// remove standard packages
	pkgs = removestdpkgs(pkgs)

	// find the unvendored packages by removing packages that contain the
	// projectImportPath as a prefix in the import path
	//
	// by using projectImportPath we also remove internal packages
	uvpkgs := removeprefix(projectImportPath, pkgs)

	// verbosity
	if verbose {
		fmt.Println(" 	" + strconv.Itoa(len(uvpkgs)) + " packages found")
	}

	// filter out vendored packages
	vpkgs := selectprefix(vendorProjectPathSlashed, pkgs)

	// check if no externally vendored or unvendored packages exist
	if len(uvpkgs) < 1 && len(vpkgs) < 1 {

		// remove everthing in the vendor directory
		os.RemoveAll(vendorPath)

		return nil
	}

	//
	// Step 3. If the vendors.yml manifest file exists, load it in memory.
	//
	////

	// create an empty slice of vendors to fill.
	var vf []vendor

	// check if vend file path exists.
	if _, err := os.Stat(vendorFilePath); err == nil {

		// verbosity
		if verbose {
			fmt.Print("loading " + vendorFilePath + "...")
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
		} else {

			if verbose {
				fmt.Println("		complete")
			}
		}

	} else {

		// verbosity
		if verbose {
			fmt.Println("			file missing: " + vendorFilePath)
		}
	}

	//
	// Step 4. Verify vendored packages and treat bad ones as unvendored packages.
	//
	////

	// check vpkgs is not empty
	if len(vpkgs) > 0 {

		// iterate over vpkgs
		for _, pkg := range vpkgs {

			// remove project path to create a complete absolute filepath
			vpath := pkg[len(vendorPath):]

			// get stats on the pkg
			if _, err := os.Stat(pkg); err != nil {

				// check if the path does not exist
				if os.IsNotExist(err) {

					// verbosity
					if verbose {
						fmt.Println("missing vendored code for " + vpath)
					}

					// clean pkg path to be unvendored
					pkg = pkg[len(vendorProjectPathSlashed):]

					// append package into the unvendored package object
					uvpkgs = append(uvpkgs, pkg)
				}

				return err
			}

			// iterate through file
			for _, v := range vf {
				fmt.Println("vendored pkgs: " + pkg + " vs " + v.Path)
			}
		}
	}

	//
	// Step 5. Identify package repositories and filter out repo subpackages.
	//
	////

	// check uvpkgs is not empty
	if len(uvpkgs) > 0 {

		// create a repo map of package paths to RepoRoots
		rmap := make(map[string]*RepoRoot)

		// verbosity
		if verbose {
			fmt.Print("identifying repositories... ")
		}

		// remove package imports that might already be included
		//
		// example: "gopkg.in/mgo.v2/bson" -> "gopkg.in/mgo.v2"
		for _, pkg := range uvpkgs {

			// iterate through file
			for _, v := range vf {

				if pkg == v.Path && len(v.Rev) > 0 {

					if _, err := os.Stat(filepath.Join(localpath, vendorProjectPath, pkg)); err == nil {
						goto UnvendoredMatch
					}
				}
			}

			// check if package path is missing from repo map
			if _, ok := rmap[pkg]; !ok {

				// determine import path dynamically by pinging repository
				r, err := RepoRootForImportPath(pkg, false)
				if err != nil {
					e := err.Error()
					msg := "no go-import meta tags"
					if e[len(e)-len(msg):] == msg {
						return errors.New("Are you behind a proxy?" + e + "\n")
					}
					return err
				}

				// add the RepoRoot to the RepoRoot map
				rmap[r.Root] = r
			}
		UnvendoredMatch:
		}

		// verbosity
		if verbose {
			fmt.Println("			complete")
		}

		//
		// Step 6. Download and vendor packages.
		//
		////

		// check that the repo map is not empty
		if len(rmap) > 0 {

			// verbosity
			if verbose {
				fmt.Println("downloading...")
			}

			// iterate through the rmap
			for _, r := range rmap {

				// create a directory for the pkg
				os.MkdirAll(filepath.Dir(filepath.Join(vendorTempPath, r.Root)), 0777)

				// iterate over the vendors in the vendor file
				for _, v := range vf {

					// check if we have a match, and a given revision exists
					if r.Root == v.Path && len(v.Rev) > 0 {

						// verbosity
						if verbose {
							fmt.Println(" ↓ " + r.Repo + " (" + v.Rev + ")")
						}

						// create the repository at that specific revision
						r.VCS.CreateAtRev(filepath.Join(vendorTempPath, r.Root), r.Repo, v.Rev)
						goto RevMatch
					}
				}

				r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo)

				// verbosity
				if verbose {
					rev, err := r.VCS.identify(filepath.Join(vendorTempPath, r.Root))
					if err != nil {
						return err
					}
					fmt.Println(" ↓ " + r.Repo + " (" + rev + ")")
				}

			RevMatch:
			}

			// verbosity
			if verbose {
				fmt.Print("vendoring packages...")
			}

			// iterate through the rmap
			for _, r := range rmap {

				// remove all
				if err := os.RemoveAll(filepath.Join(vendorPath, r.Root)); err != nil {
					return err
				}

				// mkdir
				if err := os.MkdirAll(filepath.Dir(filepath.Join(vendorPath, r.Root)), 0777); err != nil {
					return err
				}

				// copy
				if err := CopyDir(filepath.Join(vendorTempPath, r.Root), filepath.Join(vendorPath, r.Root)); err != nil {
					return err
				}
			}

			// remove the temp
			os.RemoveAll(vendorTempPath)

			// verbosity
			if verbose {
				fmt.Println("	complete")
			}
		}

		//
		// Step 7. Write the vendors.yml manifest file.
		//
		////

		//
		// Step 8. Rewrite import paths.
		//
		////
	}

	// if not in vendor file then add it to vendors
	return nil
}
