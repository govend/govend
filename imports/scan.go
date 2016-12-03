// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package imports

import (
	"go/build"
	"os"
	"sort"

	"github.com/govend/govend/imports/filters"
)

// ScanOptions represents available scan options.
type ScanOptions int

const (
	// SinglePackage only scans a single package
	SinglePackage ScanOptions = iota

	// SkipTestFiles does not scan files that end in "_test.go".
	SkipTestFiles

	// SkipFilters returns the raw unfiltered list of scanned packages.
	SkipFilters
)

// Scan takes a directory and scans it for import dependencies.
func Scan(path string, options ...ScanOptions) ([]string, error) {

	// if the path is a Godeps path, filter it out
	path = filters.Godeps([]string{path})[0]

	// parse scan options
	var singlePackage, skipTestFiles, skipFilters bool
	for _, option := range options {
		switch option {
		case SinglePackage:
			singlePackage = true
		case SkipTestFiles:
			skipTestFiles = true
		case SkipFilters:
			skipFilters = true
		}
	}

	ctx := build.Default
	mode := build.IgnoreVendor
	if singlePackage {
		mode = build.FindOnly
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	p, err := ctx.Import(path, cwd, mode)
	if err != nil {
		return nil, err
	}
	pkgs := p.Imports
	if !skipTestFiles {
		pkgs = append(pkgs, p.TestImports...)
	}

	// filter packages
	if !skipFilters {
		pkgs = filters.Exceptions(pkgs)
		pkgs = filters.Standard(pkgs)
		pkgs = filters.Local(pkgs)
		pkgs = filters.Godeps(pkgs)
	}
	pkgs = filters.Duplicates(pkgs)

	// sort
	sort.Strings(pkgs)

	return pkgs, nil
}
