package repo

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gophersaurus/govend/go15experiment"
	"github.com/gophersaurus/govend/manifest"
	"github.com/gophersaurus/govend/packages"
)

// VendCMD
func VendCMD(vendorDir, vendorFile string, verbose, recursive bool) error {

	if !go15experiment.Version() {
		return errors.New("govend only works with go version 1.5+")
	} else if !go15experiment.On() {
		return errors.New("govend currently requires the 'GO15VENDOREXPERIMENT' environment variable set to '1'")
	}

	// scan for external packages
	pkgs, err := packages.ScanDeps(".", vendorDir, false)
	if err != nil {
		return err
	}

	repomap := make(map[string]*Repo)
	for _, pkg := range pkgs {
		repo, err := Ping(pkg)
		if err != nil {
			return err
		}
		repomap[repo.ImportPath] = repo
	}

	// if the vendor manifest file exists, read it
	var vendors []manifest.Vendor
	manifestFile := filepath.Join(vendorDir, vendorFile)
	if _, err := os.Stat(manifestFile); err == nil {
		if err := manifest.Read(manifestFile, &vendors); err != nil {
			return err
		}
	}

	// the final vendors manifest file slice
	var vendorsManifest []manifest.Vendor

	// filter out vendored repositories from the repomap
	for _, vendor := range vendors {
		if _, ok := repomap[vendor.Path]; ok {
			vendorsManifest = append(vendorsManifest, vendor)
			delete(repomap, vendor.Path)
		}
	}

	// determine the absolute file path for the current local directory
	localpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// download the repository contents
	for _, repo := range repomap {
		rev, err := Download(repo, filepath.Join(localpath, vendorDir))
		if err != nil {
			return err
		}
		vendorsManifest = append(vendorsManifest, manifest.NewVendor(repo.ImportPath, rev))
	}

	// more stuff here ?

	if len(vendorsManifest) > 0 {
		if err := manifest.Write(manifestFile, &vendorsManifest); err != nil {
			return err
		}
	} else {
		os.Remove(manifestFile)
	}

	if recursive {

		// scan vendored dependencies for external packages
		rpkgs, err := packages.ScanDeps(".", vendorDir, false)
		if err != nil {
			return err
		}

		for _, pkg := range rpkgs {
			if _, err := os.Stat(filepath.Join(vendorDir, pkg)); os.IsNotExist(err) {
				if err := VendCMD(vendorDir, vendorFile, verbose, recursive); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/*

// The vend command is the default command for govend.
//
// The vend command vendors packages into the vendor directory.
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
func vendcmd(verbose, recursive bool) error {

	//
	// Step 1. Identify all relative file paths necessary for the current project.
	//

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

	//
	// Step 2. Identify all types of packages currently present in the project.
	//

	// verbosity
	if verbose {
		fmt.Print("scanning for external unvendored packages...")
	}

	// scan for external packages
	pkgs, err := scan(".", false)
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
		if err := os.RemoveAll(vendorPath); err != nil {
			return err
		}

		return nil
	}

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

		// verbosity
		if verbose {
			fmt.Println("will generate manifest... 			" + vendorFilePath)
		}
	}

	//
	// Step 4. Verify vendored packages and treat bad ones as unvendored packages.
	//

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
			if _, err := os.Stat(filepath.Join(vendorPath, vpkgpath)); os.IsNotExist(err) {

				// verbosity
				if verbose {
					fmt.Println("\n 	missing vendored code for " + vpkgpath + "...")
				}

				// append package into the unvendored package object
				uvpkgs = append(uvpkgs, vpkgpath)
			}

			// iterate through the manifest file
			for _, v := range manifest {
				if vpkgpath == v.Path {
					goto VendoredPackageMatch
				}
			}

			// add the missing vpkgpath to the new vendors manifest file
			//
			// THIS PRODUCED POOR RESULTS... WE NEED TO BE SMARTER...
			// manifest = append(manifest, vendor{Path: vpkgpath})

		VendoredPackageMatch:
		}

		// verbosity
		if verbose {
			fmt.Println("			complete")
		}
	}

	// check uvpkgs is not empty
	if len(uvpkgs) > 0 {

		//
		// Step 5. Identify package repositories and filter out repo subpackages.
		rmap, err := pingrepos(uvpkgs, manifest, localpath, vendorProjectPath, verbose)
		if err != nil {
			return err
		}

		// check that the repo map is not empty
		if len(rmap) > 0 {

			//
			// Step 6. Download and vendor packages.
			manifest, err = download(rmap, manifest, vendorTempPath, vendorPath, verbose)
			if err != nil {
				return err
			}
		}

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
*/
