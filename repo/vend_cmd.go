package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gophersaurus/govend/go15experiment"
	"github.com/gophersaurus/govend/manifest"
	"github.com/gophersaurus/govend/packages"
)

var badimports = map[string]string{}

// VendCMD
func VendCMD(verbose bool) error {

	// check that the current version of go supports vendoring
	if !go15experiment.Version() {
		return errors.New("govend only works with go version 1.5+")
	}
	if !go15experiment.On() {
		return errors.New("govend currently requires the 'GO15VENDOREXPERIMENT' environment variable set to '1'")
	}

	// load the manifest file
	m, err := manifest.Load()
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

		// use the network to gather some metadata on this repo
		repo, err := Ping(dep)
		if err != nil {
			if strings.Contains(err.Error(), "unrecognized import path") {
				badimports[dep] = "unable to ping"
				if verbose {
					fmt.Printf(" ✖ %s (bad ping)\n", dep)
				}
				continue
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
		}

		repos[repo.ImportPath] = repo
	}

	if verbose {
		fmt.Printf("%d deps scanned, %d packages skipped, %d repositories found\n", len(deps), len(badimports), len(repos))
	}

	return nil
}

/*
if verbose {
	fmt.Printf(" ↓ %s (%s)\n", repo.ImportPath, vendorRev)
}
*/

func download(dep string, m *manifest.Manifest, verbose bool) error {

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

		for _, d := range depdeps {
			if err := download(d, m, verbose); err != nil {
				return err
			}
		}
	}

	return nil
}
