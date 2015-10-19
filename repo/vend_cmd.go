package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophersaurus/govend/go15experiment"
	"github.com/gophersaurus/govend/manifest"
	"github.com/gophersaurus/govend/packages"
)

var badpkgs = map[string]string{}
var lastpkgs = []string{}

// VendCMD
func VendCMD(vendorDir, vendorFile string, verbose, recursive bool) error {

	if !go15experiment.Version() {
		return errors.New("govend only works with go version 1.5+")
	} else if !go15experiment.On() {
		return errors.New("govend currently requires the 'GO15VENDOREXPERIMENT' environment variable set to '1'")
	}

	// scan for external packages
	pkgs, invalidpkgs, err := packages.ScanDeps(".", vendorDir, false, verbose)
	if err != nil {
		return err
	}

	for pkg, msg := range invalidpkgs {
		if _, ok := badpkgs[pkg]; !ok {
			badpkgs[pkg] = msg
		}
	}

	repomap := make(map[string]*Repo)
	for _, pkg := range pkgs {

		if _, ok := badpkgs[pkg]; ok {
			continue
		}

		repo, err := Ping(pkg)
		if err != nil {
			if strings.Contains(err.Error(), "unrecognized import path") {

				badpkgs[pkg] = "unable to ping"
				if verbose {
					fmt.Printf(" ✖ %s (bad ping)\n", pkg)
				}

				continue
			}
			return err
		}

		repomap[repo.ImportPath] = repo
	}

	if verbose {
		fmt.Printf("%d packages scanned, %d packages skipped, %d repositories found\n", len(pkgs), len(badpkgs), len(repomap))
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

		vendorRev := "latest"
		for _, vendor := range vendors {
			if vendor.Path == repo.ImportPath {
				vendorRev = vendor.Rev
				break
			}
		}

		if verbose {
			fmt.Printf(" ↓ %s (%s)\n", repo.ImportPath, vendorRev)
		}

		rev, err := Download(repo, filepath.Join(localpath, vendorDir), vendorRev)
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
		rpkgs, invalidpkgs, err := packages.ScanDeps(".", vendorDir, false, verbose)
		if err != nil {
			return err
		}

		for pkg, msg := range invalidpkgs {
			if _, ok := badpkgs[pkg]; !ok {
				badpkgs[pkg] = msg
			}
		}

		fpkgs := []string{}
		for _, pkg := range rpkgs {
			if _, ok := badpkgs[pkg]; !ok {
				fpkgs = append(fpkgs, pkg)
			}
		}

		if len(lastpkgs) == len(fpkgs) {
			for i := range lastpkgs {
				if lastpkgs[i] != fpkgs[i] {
					break
				}
				if len(lastpkgs) == i+1 {
					if verbose {
						fmt.Println("it looks like some packages reference broken imports...")
					}
					return nil
				}
			}
		}

		lastpkgs = fpkgs

		for _, pkg := range fpkgs {
			if _, err := os.Stat(filepath.Join(vendorDir, pkg)); os.IsNotExist(err) {
				if verbose {
					fmt.Print("\ndownloading recursive dependencies...\n\n")
				}
				if err := VendCMD(vendorDir, vendorFile, verbose, recursive); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}
