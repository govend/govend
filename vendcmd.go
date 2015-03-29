package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jackspirou/govend/internal/_vendor/gopkg.in/yaml.v2"
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
		fmt.Println("identifying project paths... 			complete")
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
	var manifest []vendor

	// check if vend file path exists.
	if _, err := os.Stat(vendorFilePath); err == nil {

		// verbosity
		if verbose {
			fmt.Print("loading " + vendorFilePath + "...")
		}

		// read the vendors file.
		if err := load(vendorFilePath, &manifest); err != nil {
			return err
		}

		// check if the vend file is empty
		if len(manifest) < 1 {

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
			fmt.Println("will generate manifest... 			" + vendorFilePath)
		}
	}

	//
	// Step 4. Verify vendored packages and treat bad ones as unvendored packages.
	//
	////

	// check vpkgs is not empty
	if len(vpkgs) > 0 {

		// verbosity
		if verbose {
			fmt.Print("verifying vendored package...")
		}

		// iterate over vpkgs
		for _, vpkg := range vpkgs {

			// remove project path to create a complete absolute filepath
			vpkgpath := vpkg[len(vendorProjectPathSlashed):]

			// get stats on the pkg
			if _, err := os.Stat(filepath.Join(vendorPath, vpkgpath)); err != nil {

				// check if the path does not exist
				if os.IsNotExist(err) {

					// verbosity
					if verbose {
						fmt.Println("\n 	missing vendored code for " + vpkgpath + "...")
					}

					// append package into the unvendored package object
					uvpkgs = append(uvpkgs, vpkgpath)
				}

				return err
			}

			// iterate through the manifest file
			for _, v := range manifest {
				if vpkgpath == v.Path {
					goto VendoredPackageMatch
				}
			}

			// add the missing vpkgpath to the new vendors manifest file
			manifest = append(manifest, vendor{Path: vpkgpath})

		VendoredPackageMatch:
		}

		// verbosity
		if verbose {
			fmt.Println("			complete")
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

			// iterate through the manifest vendors.yml file looking for matches
			// and check if we already have vendored code and a revision number
			for _, v := range manifest {
				if pkg == v.Path {
					if _, err := os.Stat(filepath.Join(localpath, vendorProjectPath, pkg)); err == nil {
						goto UnvendoredPackageMatch
					}
				}
			}

			// check if the package is missing from RepoRoot map
			if _, ok := rmap[pkg]; !ok {

				// determine import path and repository type by asking the server
				// hosting the repository and package
				r, err := RepoRootForImportPath(pkg, false)
				if err != nil {
					e := err.Error()
					msg := "no go-import meta tags"
					if e[len(e)-len(msg):] == msg {
						return errors.New("Are you behind a proxy?" + e + "\n")
					}
					return err
				}

				// if the project package root is ont in the RepoRoot map, add it
				if _, ok := rmap[r.Root]; !ok {
					rmap[r.Root] = r
				}
			}

		UnvendoredPackageMatch:
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
				fmt.Println("downloading packages...")
			}

			// iterate through the rmap
			for key, r := range rmap {

				// create a directory for the pkg
				os.MkdirAll(filepath.Dir(filepath.Join(vendorTempPath, r.Root)), 0777)

				// iterate over the vendors in the vendor file
				for i, v := range manifest {

					// check if we have a match
					if r.Root == v.Path {

						if _, err := os.Stat(filepath.Join(vendorPath, v.Path)); err == nil {

							// verbosity
							if verbose {
								fmt.Println(" - " + r.Root + " (vendored)")
							}

							delete(rmap, key)

							goto UnvendoredManifestMatch
						}

						// check if a revision exists
						if len(v.Rev) > 0 {

							// verbosity
							if verbose {
								fmt.Println(" ↓ " + r.Repo + " (" + v.Rev + ")")
							}

							// create the repository at that specific revision
							if err := r.VCS.CreateAtRev(filepath.Join(vendorTempPath, r.Root), r.Repo, v.Rev); err != nil {
								return err
							}
						} else {

							// Create the repository from the default repository revision.
							if err := r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo); err != nil {
								return err
							}

							rev, err := r.VCS.identify(filepath.Join(vendorTempPath, r.Root))
							if err != nil {
								return err
							}

							manifest[i] = vendor{Path: v.Path, Rev: rev}

							// verbosity
							if verbose {
								fmt.Println(" ↓ " + r.Repo + " (" + rev + ")")
							}
						}
						goto UnvendoredManifestMatch
					}
				}

				// verbosity
				if verbose {
					fmt.Println(" ↓ " + r.Repo + " (latest)")
				}

				// Create the repository from the default repository revision.
				if err := r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo); err != nil {
					return err
				}

				if rev, err := r.VCS.identify(filepath.Join(vendorTempPath, r.Root)); err == nil {
					manifest = append(manifest, vendor{Path: r.Root, Rev: rev})
				} else {
					return err
				}

			UnvendoredManifestMatch:
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
				fmt.Println("				complete")
			}
		}
	}

	//
	// Step 7. Write the vendors.yml manifest file.
	//
	////

	// verbosity
	if verbose {
		fmt.Print("writing vendors.yml manifest...")
	}

	// marshal to yml
	bytes, err := yaml.Marshal(&manifest)
	if err != nil {
		return err
	}

	// write file
	ioutil.WriteFile(vendorFilePath, bytes, 0777)

	// verbosity
	if verbose {
		fmt.Println("			complete")
	}

	//
	// Step 8. Rewrite import paths.
	//
	////

	// verbosity
	if verbose {
		fmt.Print("rewriting import paths...")
	}

	// create an import replacement map to work with
	replacement := make(map[string]string)

	// fill the import replacement map
	for _, pkg := range uvpkgs {
		replacement[pkg] = filepath.Join(projectImportPath, vendorPath, pkg)
	}

	// rewrite import paths
	if err := rewrite(".", replacement); err != nil {
		return err
	}

	// verbosity
	if verbose {
		fmt.Println("			complete")
	}

	// if not in vendor file then add it to vendors
	return nil
}
