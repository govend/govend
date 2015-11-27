package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gophersaurus/govend/go15experiment"
	"github.com/gophersaurus/govend/manifest"
	"github.com/gophersaurus/govend/packages"
	"github.com/gophersaurus/govend/strutil"
)

var badimports = map[string]string{}
var m *manifest.Manifest

// VendCMD takes
func VendCMD(verbose bool) error {

	// check that the current version of go supports vendoring
	if !go15experiment.Version() {
		return errors.New("govend only works with go version 1.5+")
	}
	if !go15experiment.On() {
		return errors.New("govend currently requires the 'GO15VENDOREXPERIMENT' environment variable set to '1'")
	}

	// load the manifest file
	var err error
	m, err = manifest.Load()
	if err != nil {
		return err
	}

	// ensure the manifest file and vendor dir is in sync
	m.Sync()

	// scan the current project for external package dependencies
	deps, err := packages.ScanProject(".")
	if err != nil {
		return err
	}

	// setting some state
	badimports := map[string]string{} // bad package imports
	repos := make(map[string]*Repo)   // repositories

	// range over the external dependencies
	for _, dep := range deps {

		// check if the dep is a known bad import
		if _, ok := badimports[dep]; ok {
			continue
		}

		if err := download(dep, verbose); err != nil {
			return err
		}
	}

	if verbose {
		fmt.Printf("%d deps scanned, %d packages skipped, %d repositories found\n", len(deps), len(badimports), len(repos))
	}

	if err := m.Write(); err != nil {
		return err
	}

	return nil
}

/*
if verbose {
	fmt.Printf(" ↓ %s (%s)\n", repo.ImportPath, vendorRev)
}
*/

func download(dep string, verbose bool) error {

	// use the network to gather some metadata on this repo
	repo, err := Ping(dep)
	if err != nil {
		if strings.Contains(err.Error(), "unrecognized import path") {
			badimports[dep] = "unable to ping"
			if verbose {
				fmt.Printf(" ✖ %s (bad ping)\n", dep)
			}
		}
		return err
	}

	// check if the repo is missing from the manifest file
	if !m.Contains(repo.ImportPath) {

		if verbose {
			fmt.Printf(" ↓ %s (%s)\n", repo.ImportPath, "latest")
		}

		// download the repo
		rev, err := Download(repo, "vendor", "latest")
		if err != nil {
			return err
		}

		// append the repo the manifest file
		m.Append(manifest.NewVendor(repo.ImportPath, rev))

		depdeps, err := packages.Scan(filepath.Join("vendor", dep))
		if err != nil {
			return err
		}

		depdeps = packages.FilterStdPkgs(depdeps)

		projectpath, err := packages.ImportPath(filepath.Join("vendor", dep))
		if err != nil {
			return err
		}

		// filter out packages internal to the project
		depdeps = strutil.RemovePrefixInStringSlice(projectpath, depdeps)

		for _, d := range depdeps {
			if err := download(d, verbose); err != nil {
				return err
			}
		}
	}

	return nil
}
