package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gophersaurus/govend/go15experiment"
	"github.com/gophersaurus/govend/manifest"
)

// InstallCMD
func InstallCMD(vendorDir, vendorFile string, verbose bool) error {

	if !go15experiment.Version() {
		return errors.New("govend only works with go version 1.5+")
	} else if !go15experiment.On() {
		return errors.New("govend currently requires the 'GO15VENDOREXPERIMENT' environment variable set to '1'")
	}

	// if the vendor manifest file exists, read it
	var vendors []manifest.Vendor
	if _, err := os.Stat(vendorFile); err != nil {
		if os.IsNotExist(err) {
			return errors.New("unable to find/read vendor the manifest file...")
		}
		return err
	}

	if err := manifest.Read(vendorFile, &vendors); err != nil {
		return err
	}

	if verbose {
		fmt.Println("reset vendor directory...")
	}

	// remove the vendor directory
	os.RemoveAll(vendorDir)

	// determine the absolute file path for the current local directory
	localpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// download the repos specified in the vendor manifest
	for _, vendor := range vendors {
		if verbose {
			fmt.Printf(" ↓ %s (%s)\n", vendor.Path, vendor.Rev)
		}
		repo, err := Ping(vendor.Path)
		if err != nil {
			return err
		}
		_, err = Download(repo, filepath.Join(localpath, vendorDir), vendor.Rev)
		if err != nil {
			return err
		}
	}

	return nil
}
