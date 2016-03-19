// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package imports

import (
	"os"
	"sort"
	"strings"

	"github.com/govend/govend/imports/filters"
	"github.com/kr/fs"
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

	// get directory info
	dinfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// if we are in a directory step to it!
	w := fs.Walk(path)
	if dinfo.IsDir() {
		w.Step()
	}

	// we will parse a list of packages
	pkgs := []string{}
	for w.Step() {

		// skip all files and directories that start with '.' or '_'
		finfo := w.Stat()

		// check for any walker errors
		if w.Err() != nil {
			return nil, w.Err()
		}

		firstchar := []rune(finfo.Name())[0]
		if firstchar == '_' || firstchar == '.' {
			if finfo.IsDir() {
				w.SkipDir()
				continue
			} else {
				continue
			}
		}

		// skip directories named "vendor"
		if finfo.IsDir() {
			if finfo.Name() == "vendor" || finfo.Name() == "Godeps" {
				w.SkipDir()
				continue
			}
			if singlePackage {
				w.SkipDir()
				continue
			}
		}

		// if testfiles is false then skip all go tests deps
		if skipTestFiles && strings.HasSuffix(finfo.Name(), "_test.go") {
			continue
		}

		// only parse .go files
		fpath := w.Path()
		if strings.HasSuffix(fpath, ".go") {
			p, err := Parse(w.Path())
			if err != nil {

				// if the error is because of a bad file, skip the file
				if strings.Contains(err.Error(), eofError) {
					continue
				}
			}
			pkgs = append(pkgs, p...)
		}
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
