package govend

// The update command.
//
// The update command updates vendored packages.
// It achieves this by executing the following steps:
//
// Step 1. Identify all relative file paths necessary for the current project.
// Step 2. Identify all types of packages currently present in the project.
// Step 3. If the vendors.yml manifest file exists, load it in memory.
// Step 4. Check the package specified exists in the project and in the vendors manifest.
// Step 5. Check the package revision is different.
// Step 6. Download and vendor the package.
// Step 7. Write the vendors.yml manifest file.
// Step 8. Rewrite import paths.
//

/*


func updatecmd(pkg, rev string, verbose, recursive bool) error {

	if len(pkg) < 1 {
		panic("missing pkg to update")
	}

	if len(rev) > 1 {
		panic("updating to a particular revision is not current available")
	}

	//
	// Step 1. Identify all relative file paths necessary for the current project.
	//


*/

// determine the absolute file path for the current local directory
/*
	localpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
*/

/*



	// determine the project import path
	projectImportPath, err := importpath(".")
	if err != nil {
		return err
	}

	// define vendor directory paths
	vendorProjectPath := filepath.Join(projectImportPath, vendorPath)
	vendorProjectPathSlashed := vendorProjectPath + string(filepath.Separator)

	//
	// Step 2. Identify all types of packages currently present in the project.
	//

	// verbosity
	if verbose {
		fmt.Print("scanning vendored packages...")
	}

	// scan for external packages
	pkgs, err := scan(".", false)
	if err != nil {
		return err
	}

	// remove standard packages
	pkgs = removestdpkgs(pkgs)

	// filter out vendored packages
	vpkgs := selectprefix(vendorProjectPathSlashed, pkgs)

	//
	// Step 3. If the vendors.yml manifest file exists, load it in memory.
	//

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
		panic("missing vendor manifest...")
	}

	//
	// Step 4. Verify vendored packages and treat bad ones as unvendored packages.
	// Step 4. Check the package specified exists in the project and in the vendors manifest.
	//

	// check vpkgs is not empty
	if len(vpkgs) < 1 {
		panic(pkg + " is not present...")
	}

	// iterate over vpkgs
	for _, vpkg := range vpkgs {

		// remove project path to create a complete absolute filepath
		vpkgpath := vpkg[len(vendorProjectPathSlashed):]

		if vpkgpath == pkg {

			// iterate through the manifest file
			for i, v := range manifest {
				if pkg == v.Path {

					r, err := RepoRootForImportPath(pkg, false)
					if err != nil {
						e := err.Error()
						msg := "no go-import meta tags"
						if e[len(e)-len(msg):] == msg {
							return errors.New("Are you behind a proxy?" + e + "\n")
						}
						return err
					}

					// create a temp directory for the pkg
					if err := os.MkdirAll(filepath.Dir(filepath.Join(vendorTempPath, r.Root)), 0777); err != nil {
						return err
					}

					// verbosity
					if verbose {
						fmt.Println(" ↓ " + r.Repo + " (latest)")
					}

					// Create the repository from the default repository revision.
					if err := r.VCS.Create(filepath.Join(vendorTempPath, r.Root), r.Repo); err != nil {
						return err
					}

					rev, err := r.VCS.identify(filepath.Join(vendorTempPath, r.Root))

					if err != nil {
						return err
					}

					manifest[i] = vendor{Path: v.Path, Rev: rev}

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

					// remove the temp
					os.RemoveAll(vendorTempPath)

					// verbosity
					if verbose {
						fmt.Println("				complete")
					}

					//
					// Step 7. Write the vendors.yml manifest file.
					//

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

					// verbosity
					if verbose {
						fmt.Print("rewriting import paths...")
					}

					// create an import replacement map to work with
					replacement := make(map[string]string)

					// fill the import replacement map
					replacement[pkg] = filepath.Join(projectImportPath, vendorPath, pkg)

					// rewrite import paths
					if err := rewrite(".", replacement); err != nil {
						return err
					}

					// verbosity
					if verbose {
						fmt.Println("			complete")
					}

					if recursive {

						// scan for external packages
						tpkgs, err := scan(".", false)
						if err != nil {
							return err
						}

						// remove standard packages
						tpkgs = removestdpkgs(tpkgs)

						// find the unvendored packages by removing packages that contain the
						// projectImportPath as a prefix in the import path
						//
						// by using projectImportPath we also remove internal packages
						tuvpkgs := removeprefix(projectImportPath, tpkgs)

						if len(tuvpkgs) > 0 {

							if verbose {
								fmt.Println("")
								fmt.Println("recursively vendoring dependencies...")
								fmt.Println("")
							}
							if err := vendcmd(verbose, recursive); err != nil {
								return err
							}
						}
					}

					return nil

				}
			}
		}
	}

	panic(pkg + " is not present...")

}
*/
