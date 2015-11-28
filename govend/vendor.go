package govend

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophersaurus/govend/go15experiment"
	"github.com/gophersaurus/govend/manifest"
	"github.com/gophersaurus/govend/packages"
	"github.com/gophersaurus/govend/repo"
)

var m *manifest.Manifest
var badimports map[string]string
var lastdep string

// Vendor takes
func Vendor(pkgs []string, update, verbose, commands bool, format string) error {

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
	badimports = map[string]string{} // bad package imports

	// range over the external dependencies
	for _, dep := range deps {

		// check if the dep is a known bad import
		if _, ok := badimports[dep]; ok {
			continue
		}

		// download that dependency and any external deps it has
		lastdep = dep
		if err := downloadDeps(dep, verbose); err != nil {
			return err
		}
	}

	if verbose {
		fmt.Printf("%d packages scanned, %d skipped\n", len(deps), len(badimports))
	}

	if err := m.Write(); err != nil {
		return err
	}

	return nil
}

func downloadDeps(dep string, verbose bool) error {

	// use the network to gather some metadata on this repo
	r, err := repo.Ping(dep)
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
	if !m.Contains(r.ImportPath) {

		if verbose {
			fmt.Printf(" ↓ %s (%s)\n", r.ImportPath, "latest")
		}

		// download the repo
		rev, err := repo.Download(r, "vendor", "latest")
		if err != nil {
			return err
		}

		// append the repo to the manifest file
		m.Append(manifest.NewVendor(r.ImportPath, rev))
	}

	depdeps, err := packages.Scan(filepath.Join("vendor", dep))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	depdeps = packages.FilterStdPkgs(depdeps)

	for _, d := range depdeps {
		if d != dep && d != lastdep {
			lastdep = dep
			if err := downloadDeps(d, verbose); err != nil {
				return err
			}
		}
	}

	return nil
}
