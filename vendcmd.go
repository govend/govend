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

	vendorRoot := filepath.Join(projectpath, vendorPath)
	vendorRootSlash := vendorRoot + string(filepath.Separator)

	// find the unvendored packages by removing packages that contain the
	// projectpath as a prefix in the import path
	//
	// by using projectpath we also remove internal packages
	uvpkgs := removeprefix(projectpath, pkgs)

	// verbosity
	if verbose {
		fmt.Println(" 	" + strconv.Itoa(len(uvpkgs)) + " packages found")
		for _, pkg := range uvpkgs {
			fmt.Println("	" + pkg)
		}
	}

	// filter out vendored packages
	vpkgs := selectprefix(vendorRootSlash, pkgs)

	// check if no externally vendored or unvendored packages exist
	if len(uvpkgs) < 1 && len(vpkgs) < 1 {

		// remove everthing in the vendor directory
		os.RemoveAll(vendorPath)

		return nil
	}

	// check vpkgs is not empty
	if len(vpkgs) > 0 {

		// iterate over vpkgs
		for _, pkg := range vpkgs {

			// remove project path to create a complete absolute filepath
			vpath := pkg[len(vendorPath):]

			// get stats on the pkg
			if _, err := os.Stat(filepath.Join(localpath, vpath)); err != nil {

				// check if the path does not exist
				if os.IsNotExist(err) {

					// verbosity
					if verbose {
						fmt.Println("missing vendored code for " + vpath)
					}

					// clean pkg path to be unvendored
					pkg = pkg[len(vendorRootSlash):]

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

	// check uvpkgs is not empty
	if len(uvpkgs) > 0 {

		// create a repo map of package paths to RepoRoots
		rmap := make(map[string]*vcs.RepoRoot)

		// iterate over uvpkgs
		// remove package imports that might already be included
		// example: "gopkg.in/mgo.v2/bson" -> "gopkg.in/mgo.v2"
		for _, pkg := range uvpkgs {

			if verbose {
				fmt.Print("pinging... ")
			}

			// determine import path dynamically by pinging repository
			r, err := vcs.RepoRootForImportDynamic(pkg, false)
			if err != nil {
				e := err.Error()
				fmt.Println(e)
				msg := "no go-import meta tags"
				if e[len(e)-len(msg):] == msg {
					return errors.New("Are you behind a proxy?" + e + "\n")
				}
				return err
			}

			// check if package path is missing from repo map
			if _, ok := rmap[pkg]; !ok {

				// add the RepoRoot to the repo map
				rmap[pkg] = r
			}

			if verbose {
				fmt.Println("					" + pkg)
			}
		}

		// check that the repo map is not empty
		if len(rmap) > 0 {

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
							fmt.Print("downloading...")
							fmt.Println("					" + r.Repo + " " + v.Rev)
						}

						// create the repository at that specific revision
						r.VCS.CreateAtRev(filepath.Join(vendorTempPath, r.Root), r.Repo, v.Rev)
						goto RevMatch
					}
				}

				r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo)

				// verbosity
				if verbose {
					fmt.Print("downloading...")
					fmt.Println("					" + r.Repo)
				}

			RevMatch:
			}

			// iterate through the rmap
			for _, r := range rmap {

				// verbosity
				if verbose {
					fmt.Print("cleaning, ")
				}
				os.RemoveAll(filepath.Join(vendorPath, r.Root))

				// verbosity
				if verbose {
					fmt.Print("creating, ")
				}
				os.MkdirAll(filepath.Dir(filepath.Join(vendorPath, r.Root)), 0777)

				// verbosity
				if verbose {
					fmt.Print("and copying...")
				}
				CopyDir(filepath.Join(vendorTempPath, r.Root), filepath.Join(vendorPath, r.Root))

				// verbosity
				if verbose {
					fmt.Println("		" + r.Root)
				}
			}

			// verbosity
			if verbose {
				fmt.Print("removing all temporary directories...")
			}
			os.RemoveAll(vendorTempPath)

			// verbosity
			if verbose {
				fmt.Println("		complete")
			}
		}
	}

	// if not in vendor file then add it to vendors
	return nil
}
