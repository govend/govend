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

// Vend is the main function govend uses to vendor external packages.
func Vend(pkgs []string, update, verbose, tree, results, lock, hold bool, format string) error {

	// check the site is vendorable
	if err := Vendorable(); err != nil {
		return err
	}

	// attempt to load the manifest file
	m, err := manifest.Load(format)
	if err != nil {
		return err
	}

	// sync ensures that if a vendor is specified in the manifest, that the
	// repository structure is also currently present in the vendor directory,
	// this allows us to trust the manifest file
	mlen := 0
	if lock || hold || update {
		// it is important to save the manifest length before syncing, so that
		// we can tell the difference and update the manifest file
		mlen = m.Len()
		m.Sync()
	}

	// if no packages were provided as arguments, assume the current directory is
	// a go project and scan it for external packages.
	if len(pkgs) == 0 {
		pkgs, err = imports.Scan(".")
		if err != nil {
			return err
		}
	}

	// download that dependency and any external deps it has
	pkglist := map[string]bool{}
	pkgs = filters.Ellipses(pkgs)
	stack := newVendorStack(pkgs...)
	for !stack.empty() {

		// pop an import package path off the stack
		pkg := stack.pop()
		if _, ok := pkglist[pkg.path]; ok {
			continue
		}

		// use the network to gather some metadata on this repo
		repo, err := repos.Ping(pkg.path)
		if err != nil {
			fmt.Printf("%s (bad ping): %s\n", pkg.path, err)
			pkglist[pkg.path] = false
			continue
		}

		if _, ok := pkglist[repo.ImportPath]; ok {
			continue
		}

		if verbose {
			if tree {
				writeBlanks(pkg.level)
			}
			fmt.Printf("%s\n", repo.ImportPath)
		}

		// check if the repo is missing from the manifest file
		vpath := filepath.Join("vendor", repo.ImportPath)
		if !m.Contains(repo.ImportPath) && !dirExists(vpath) || lock || hold || update {
			rev, err := repos.Download(repo, "vendor", "latest")
			if err != nil {
				fmt.Printf("%s (download error): %s\n", repo.ImportPath, err)
				pkglist[pkg.path] = false
				continue
			}
			m.Append(repo.ImportPath, rev, hold)
		} else {
			for _, vendor := range m.Vendors {
				if vendor.Path == repo.ImportPath {
					if _, err := repos.Download(repo, "vendor", vendor.Rev); err != nil {
						fmt.Printf("%s (download error): %s\n", repo.ImportPath, err)
						pkglist[pkg.path] = false
						continue
					}
				}
			}
		}

		vpkg := filepath.Join("vendor", pkg.path)
		importOptions := []imports.ScanOptions{}
		if !hold {
			importOptions = append(importOptions, imports.SinglePackage)
		}
		deps, err := imports.Scan(vpkg, importOptions...)
		if err != nil {
			fmt.Printf("%s (scan error): %s\n", pkg.path, err)
			continue
		}
		pkglist[pkg.path] = true
		pkglist[repo.ImportPath] = true

		// push
		if len(deps) > 0 {
			stack.push(pkg.level+1, deps...)
		}
	}

	if verbose && results {
		fmt.Printf("\npackages scanned: %d\n", len(pkglist))
		skipped := 0
		for _, ok := range pkglist {
			if !ok {
				skipped++
			}
		}
		fmt.Printf("packages skipped: %d\n", skipped)
		fmt.Printf("repos downloaded: %d\n", m.Len())
	}

	if lock || hold || mlen > 0 {
		if err := m.Write(); err != nil {
			return err
		}
	}

	return nil
}

// writeBlanks writes a number of blank spaces.
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

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
