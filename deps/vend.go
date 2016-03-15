// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/govend/govend/deps/repos"
	"github.com/govend/govend/imports"
	"github.com/govend/govend/imports/filters"
	"github.com/govend/govend/manifest"
)

// VendOptions represents available vend options.
type VendOptions int

const (
	// Update updates vendored repositories.
	Update VendOptions = iota

	// Lock locks the revision version of vendored repositories.
	Lock

	// Hold holds onto a vendored repository, even if none of its import paths
	// are used in the project source code.
	Hold

	// Prune removes vendored packages that are not needed.
	Prune

	// Verbose prints out packages as they are vendored.
	Verbose

	// Tree prints the names of packages as an indented tree.
	Tree

	// Results prints a summary of the number of packages scanned, packages
	// skipped, and repositories downloaded.
	Results
)

// Vend is the main function govend uses to vendor external packages.
func Vend(pkgs []string, format string, options ...VendOptions) error {
	// parse vend options
	var update, lock, hold, prune, verbose, tree, results bool
	for _, option := range options {
		switch option {
		case Update:
			update = true
		case Lock:
			lock = true
		case Hold:
			hold = true
		case Prune:
			prune = true
		case Verbose:
			verbose = true
		case Tree:
			tree = true
		case Results:
			results = true
		}
	}

	// load the manifest file if it exists
	m, err := manifest.Load(format)
	if err != nil {
		return err
	}

	// only sync the manfiest file for the lock, hold, or update flags
	if lock || hold || update {
		// sync ensures that if a vendored repository is specified in the manifest
		// file, the same repository directory structure also exists inside the
		// vendor directory
		m.Sync()
	}

	// check if any import package paths were provided as arguments
	if len(pkgs) == 0 {
		// if no packages were provided, we can only assume the
		// current directory contains go source code, so we scan it
		pkgs, err = imports.Scan(".")
		if err != nil {
			return err
		}
	}

	// rather than using some variable state to track if a repository has been
	// downloaded we check if that repo or package import path exists inside the
	// vendor directory structure
	//
	// inorder to trust the vendor directory structure for lock, hold or update
	// flags we must remove the contents of the vendor directory before vendoring
	if lock || update {
		if err := os.RemoveAll("vendor"); err != nil {
			return err
		}
	}

	// support the go get "/.../" and "/..." ellipses syntax by filtering it out
	pkgs = filters.Ellipses(pkgs)

	// cache is a cache map of package import paths to a boolean results value
	cache := map[string]bool{}

	// newVendorStack takes a list of packages, reverses them, and places them
	// on a stack of a slice of strings
	//
	// the stack data structure allows for the tree flag implimentation to be
	// clean and simple
	stack := newVendorStack(pkgs...)
	keepers := pkgs

	// iterating over a stack allows vendoring of packages and any of their
	// dependencies to the nth degree in the order they are discovered
	for !stack.empty() {
		pkg := stack.pop()

		// if the package path has been cached, skip it despite the result value
		if _, ok := cache[pkg.path]; ok {
			continue
		}

		// tell the humans we are going to process this package
		if verbose {
			if tree {
				writeBlanks(pkg.level)
			}
			fmt.Printf("%s\n", pkg.path)
		}

		// if the package import path structure does not exist inside the vendor,
		// we need to download it
		if !vendorDirExists(pkg.path) {

			// ping the VCS repo across the network to gather metadata tags using
			// the package import path
			repo, err := repos.Ping(pkg.path)
			if err != nil {
				fmt.Printf("%s (bad ping): %s\n", pkg.path, err)
				cache[pkg.path] = false
				continue
			}

			// if the manifest file contains the repo and the update flag is off then
			// use the manifest revision version even if the value is an empty string
			//
			// otherwise get the latest version of the repository
			target := "latest"
			if vendor, ok := m.Contains(repo.ImportPath); ok && !update {
				target = vendor.Rev
			}

			rev, err := repos.Download(repo, "vendor", target)
			if err != nil {
				fmt.Printf("%s (download error): %s\n", repo.ImportPath, err)
				cache[pkg.path] = false
				continue
			}
			m.Append(repo.ImportPath, rev, hold)
		}

		// update the cache for the package import path after checking if the
		// package existed inside the vendor directory, this allows for both
		// the download neede and download not needed use cases
		cache[pkg.path] = true

		// we must scan the recently vendored package for any dependencies it
		// relies on so we can vendor them on the next iteration
		//
		// if the hold flag is applied we must vendor all the packages in that repo
		// since they might be used for tooling and not source code dependencies
		vdeps := []string{}
		if hold {
			vdeps, err = imports.Scan(filepath.Join("vendor", pkg.path))
		} else {
			vpath := filepath.Join("vendor", pkg.path)
			vdeps, err = imports.Scan(vpath, imports.SinglePackage)
		}
		if err != nil {
			fmt.Printf("%s (scan error): %s\n", pkg.path, err)
			continue
		}

		// push the vendor package dependencies on the stack so we can parse sub
		// packages in the order they are discovered
		if len(vdeps) > 0 {
			stack.push(pkg.level+1, vdeps...)
			keepers = append(keepers, vdeps...)
		}
	}

	// tell the humans the results
	if verbose && results {
		fmt.Printf("\npackages scanned: %d\n", len(cache))
		skipped := 0
		for _, ok := range cache {
			if !ok {
				skipped++
			}
		}
		fmt.Printf("packages skipped: %d\n", skipped)
		fmt.Printf("repos downloaded: %d\n", m.Len())
	}

	// if we need to do so, update the manifest file
	if lock || hold || update && fileExists(m.Filename()) {
		if err := m.Write(); err != nil {
			return err
		}
	}

	if prune {
		if verbose {
			fmt.Print("\npruning packages... ")
		}
		keepers = filters.Duplicates(keepers)
		prunePackages(keepers)
		if verbose {
			fmt.Println("finished")
		}
	}

	return nil
}

func writeBlanks(num int) {
	num = num * 2
	for num > 0 {
		fmt.Printf(" ")
		num--
	}
}

func badImportPath(err error) bool {
	return strings.Contains(err.Error(), "unrecognized import path")
}

func vendorDirExists(path string) bool {
	_, err := os.Stat(filepath.Join("vendor", path))
	return !os.IsNotExist(err)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
