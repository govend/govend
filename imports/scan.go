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

type ScanOption func(*scanOptions)

var (
	// AllBuildTags includes files of every build tag.
	AllBuildTags ScanOption

	// BuildTags sets extra build tags to be included when considering which
	// package imports will be needed. It defaults to go/build's default
	// Context tags, which is to say only OS/Arch related tags.
	BuildTags func(...string) ScanOption

	// SkipFilters returns the raw unfiltered list of scanned packages.
	SkipFilters ScanOption

	// SkipTestFiles does not scan files that end in "_test.go".
	SkipTestFiles ScanOption
)

type scanOptions struct {
	allTags       bool
	buildTags     []string
	skipFilters   bool
	skipTestFiles bool
}

func init() {
	AllBuildTags = func(opts *scanOptions) { opts.allTags = true }
	BuildTags = func(tags ...string) ScanOption {
		return func(opts *scanOptions) {
			opts.buildTags = append(opts.buildTags, tags...)
		}
	}
	SkipFilters = func(opts *scanOptions) { opts.skipFilters = true }
	SkipTestFiles = func(opts *scanOptions) { opts.skipTestFiles = true }
}

// Scan takes a directory and scans it for import dependencies.
func Scan(path string, options ...ScanOption) ([]string, error) {

	// if the path is a Godeps path, filter it out
	path = filters.Godeps([]string{path})[0]

	// parse scan options
	var opts scanOptions
	for _, opt := range options {
		opt(&opts)
	}

	// init build context for scanning imports
	ctx := build.Default
	if len(opts.buildTags) > 0 {
		ctx.BuildTags = append(ctx.BuildTags, opts.buildTags...)
	}
	if opts.allTags {
		ctx.UseAllFiles = true
	}

	// Find imports
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	p, err := ctx.Import(path, cwd, build.IgnoreVendor)
	if err != nil {
		return nil, err
	}
	pkgs := p.Imports
	if !opts.skipTestFiles {
		pkgs = append(pkgs, p.TestImports...)
	}

	// filter packages
	if !opts.skipFilters {
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
