// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

// Package deps provides vendoring for repositories.
package deps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/govend/govend/deps/repos"
	"github.com/govend/govend/imports"
	"github.com/govend/govend/imports/filters"
	"github.com/govend/govend/manifest"
)

// Vend is the main function govend uses to vendor external packages.
// Vend also invokes the Hold and Prune methods.
func Vend(pkgs []string, format string, options ...VendOptions) error {

	// parse VendOptions into usable boolean values
	update, lock, hold, prune, ignore, verbose, tree, results := parseVendOptions(options)

	// load or create an empty manifest file
	m, err := manifest.Load(format)
	if err != nil {
		return err
	}

	// sync ensures that if a vendored repository is specified in the manifest
	// file, the same repository directory structure also exists inside the
	// vendor directory
	if lock || hold || update {
		m.Sync()
	}

	// if no packages were provided, we can only assume the current relative
	// directory contains Go source code, therefore so we should scan it
	if len(pkgs) == 0 && !ignore {
		pkgs, err = imports.Scan(".")
		if err != nil {
			return err
		}
	} else {
		for _, pkg := range m.Vendors {
			pkgs = append(pkgs, pkg.Path)
		}
	}

	// to be efficent, we need to track if a repository has already been
	// downloaded
	//
	// pinging and downloading a repository twice is a waste of valuable network
	// bandwidth, time and disk usage
	//
	// rather than using state to track if a repository has been download, we
	// can instead check if the repository import path exists inside of the
	// vendor directory
	//
	// to trust the contents of the vendor directory, we must first remove it
	// before starting the vendoring process
	//
	// removing the vendor directory ensures the repository package paths found
	// do not originate from a previous vendoring session
	//
	// removing the vendor directory also helps keeps the vendored repositories
	// clean and fresh, preventing stale packages from sticking around when they
	// are no longer needed
	if lock || update {
		if err := os.RemoveAll("vendor"); err != nil {
			return err
		}
	}

	// the go get "/.../" and "/..." ellipses syntax should be supported since
	// govend is based on the go get command
	//
	// the purpose of the ellipses syntax is to iterate through all dependent
	// packages and repositories to the nth degree
	//
	// govend iterates through packages to the nth degree by default since it is
	// a vendoring tool, so ellipses syntax should simply be filtered it out
	pkgs = filters.Ellipses(pkgs)

	// the stack data structure allows for the dependency tree printing
	// implementation to be simple and clean
	//
	// newStack reverses a list of packages and places them on the stack
	stack := newStack(pkgs...)

	// we need some state to track the dependency tree as well as a cache of
	// parsed packages
	//
	// deptree is a list of package import paths that describe the valid
	// dependency tree, it is used for pruning
	//
	// cache is a cache map of package import paths to a boolean results value
	deptree := pkgs
	cache := map[string]bool{}

	// iterating over a stack allows vendoring of packages and any of their
	// dependencies to the nth degree in the order they are discovered
	for !stack.empty() {

		// pop the next package path off the stack
		pkg := stack.pop()

		// skip package paths that have already been cached
		if _, ok := cache[pkg.path]; ok {
			continue
		}

		// write level * 2 blanks to visualize the package in the dep tree
		if verbose && tree {
			writeDoubleBlanks(pkg.level)
		}

		// print the package import path relevant to the $GOPATH/src
		if verbose {
			fmt.Printf("%s\n", pkg.path)
		}

		// check if the import package path exists inside the vendor directory
		if _, err := os.Stat(filepath.Join("vendor", pkg.path)); os.IsNotExist(err) {

			// we know the package import path does not exist inside of the vendor
			// directory, but we don't know if the import path is representative of
			// the repository url
			//
			// we need to get info on the VCS repo which contains this package by
			// "pinging" it across the network, thereby gathering metadata tags
			// provided/exposed by VCS server host
			repo, err := repos.Ping(pkg.path)
			if err != nil {
				reportBadPing(pkg.path, err)
				cache[pkg.path] = false
				continue
			}

			// if the manifest file contains the repo and the update flag is off then
			// use the manifest revision version even if the value is an empty string
			//
			// otherwise get the latest version of the repository
			revision := "latest"
			if vendor, ok := m.Contains(repo.ImportPath); ok && !update {
				revision = vendor.Rev
			}

			// download the repository at the requested target revision into the
			// vendor directory, the revision returned is the actual one downloaded
			rev, err := repos.Download(repo, "vendor", revision)
			if err != nil {
				reportBadPing(repo.ImportPath, err)
				cache[pkg.path] = false
				continue
			}

			// if the repository path already exists in the manifest file appending
			// does not add a duplicate, it simply overwrites the current values
			m.Append(repo.ImportPath, rev, hold)
		}

		// update the cache to include the package import path
		cache[pkg.path] = true

		// we need to scan the recently vendored package for any dependencies that
		// it relies on so that they can be vendored in the next iterations
		//
		// but... first we need to determine which scan options provided
		scanOpts := []imports.ScanOptions{}
		if !hold {
			scanOpts = append(scanOpts, imports.SinglePackage)
			if prune {
				scanOpts = append(scanOpts, imports.SkipTestFiles)
			}
		}
		vdeps, err := imports.Scan(filepath.Join("vendor", pkg.path), scanOpts...)
		if err != nil {
			reportBadPing(pkg.path, err)
			continue
		}

		// push the vendor package dependencies on the stack so we can parse
		// sub/dependent packages in the order they are discovered
		//
		// also add the vendor package dependencies to the deptree for the pruning
		// process later, if we need it
		if len(vdeps) > 0 {
			stack.push(pkg.level+1, vdeps...)
			deptree = append(deptree, vdeps...)
		}
	}

	// download any repositories that are on hold
	numOfReposOnHold := Hold(m, verbose)

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
		fmt.Printf("repos downloaded: %d\n", m.Len()+numOfReposOnHold)
		if numOfReposOnHold > 0 {
			fmt.Printf("repos being held: %d\n", numOfReposOnHold)
		}
	}

	// if a lock or hold flag is present, or if and update was requested and a
	// manifest file currently exists on disk, then update the manifest file
	if lock || hold || update && fileExists(m.Filename()) {
		if err := m.Write(); err != nil {
			return err
		}
	}

	if prune {
		Prune(deptree, verbose)
	}

	return nil
}

func reportBadPing(path string, err error) {
	fmt.Printf("%s bad ping: %s\n", path, err)
}

func writeDoubleBlanks(num int) {
	num = num * 2
	for num > 0 {
		fmt.Printf(" ")
		num--
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
