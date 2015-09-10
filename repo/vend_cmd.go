package repo

import (
	"errors"
	"fmt"
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

	if verbose {
		fmt.Printf("%d packages scanned, %d repositories found\n", len(pkgs), len(repomap))
	}

	// if the vendor manifest file exists, read it
	var vendors []manifest.Vendor
	if _, err := os.Stat(vendorFile); err == nil {
		if err := manifest.Read(vendorFile, &vendors); err != nil {
			return err
		}
	}

	// the final vendors manifest file slice
	var vendorsManifest []manifest.Vendor

	// filter out vendored repositories from the repomap
	for _, vendor := range vendors {
		if _, ok := repomap[vendor.Path]; ok {
			if _, err := os.Stat(filepath.Join(vendorDir, vendor.Path)); err == nil {
				delete(repomap, vendor.Path)
			}
			vendorsManifest = append(vendorsManifest, vendor)
		}
	}

	// determine the absolute file path for the current local directory
	localpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// download the repository contents
	for _, repo := range repomap {
		revision := "latest"
		for _, vendor := range vendors {
			if vendor.Path == repo.ImportPath {
				revision = vendor.Rev
				break
			}
		}
		if verbose {
			fmt.Printf(" ↓ %s (%s)\n", repo.ImportPath, revision)
		}
		rev, err := Download(repo, filepath.Join(localpath, vendorDir), revision)
		if err != nil {
			return err
		}
		vendorsManifest = append(vendorsManifest, manifest.NewVendor(repo.ImportPath, rev))
	}

	if len(vendorsManifest) > 0 {
		if err := manifest.Write(vendorFile, &vendorsManifest); err != nil {
			return err
		}
	} else {
		os.Remove(vendorFile)
	}

	if recursive {

		// scan vendored dependencies for external packages
		rpkgs, err := packages.ScanDeps(".", vendorDir, false)
		if err != nil {
			return err
		}

		for _, pkg := range rpkgs {
			if _, err := os.Stat(filepath.Join(vendorDir, pkg)); os.IsNotExist(err) {
				if verbose {
					fmt.Print("\ndownloading recursive dependencies...\n\n")
				}
				if err := VendCMD(vendorDir, vendorFile, verbose, recursive); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
