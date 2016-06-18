//go:generate go run generate.go

// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

// Package filters provides filters for Go package import paths.
package filters

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Standard filters out standard packages.
func Standard(pkgs []string) []string {

	var list []string
	for _, pkg := range pkgs {

		// skip CGO and any relative import paths
		if pkg == "C" || pkg[0] == '.' {
			continue
		}

		name := path.Base(pkg)
		if stds, ok := stdpkgs[name]; ok {
			for _, stdpkg := range stds {
				if pkg == stdpkg.path {
					goto SKIP
				}
			}
		}
		list = append(list, pkg)

	SKIP:
	}

	return list
}

// Exceptions filters out very special exceptional packages.
func Exceptions(pkgs []string) []string {

	var list []string
	for _, pkg := range pkgs {

		for _, exception := range ExceptionList {
			if Match(pkg, exception) {
				goto SKIP
			}
		}
		list = append(list, pkg)

	SKIP:
	}

	return list
}

// Duplicates filters out any duplicate packages.
func Duplicates(pkgs []string) []string {
	result := []string{}
	seen := map[string]string{}
	for _, val := range pkgs {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

// Local filters out any local packages.
func Local(pkgs []string) []string {
	projectPath := projectImportPath()
	prefix := projectPath + "/"

	var list []string
	for _, pkg := range pkgs {
		if pkg == projectPath || strings.HasPrefix(pkg, prefix) {
			continue
		}

		list = append(list, pkg)
	}

	return list
}

// Ellipses trims the ellipses suffix off of package import paths.
func Ellipses(pkgs []string) []string {
	var list []string
	for _, pkg := range pkgs {
		pkg = strings.TrimSuffix(pkg, "/.../")
		pkg = strings.TrimSuffix(pkg, "/...")
		list = append(list, pkg)
	}
	return list
}

// Godeps filters out Godeps package import paths.
func Godeps(pkgs []string) []string {
	var list []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, "/Godeps/_workspace/src/") {
			split := strings.SplitAfter(pkg, "/Godeps/_workspace/src/")
			if len(split) > 1 {
				pkg = split[1]
			}
		}
		list = append(list, pkg)
	}
	return list
}

// projectImportPath returns the import path of the current working directory.
// It does so by trimming off the $GOPATH/src prefix.
func projectImportPath() string {
	// determine the current working directory and coerce it to an absolute
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cwd, err = filepath.Abs(cwd)
	if err != nil {
		log.Fatal(err)
	}

	// check the env $GOPATH is valid
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		log.Fatal(errors.New("$GOPATH not set"))
	}

	// trim the $GOPATH/src off the current working directory
	gosrc := filepath.Join(gopath, "src") + string(filepath.Separator)
	importpath := strings.TrimPrefix(cwd, gosrc)
	dir, file := filepath.Split(importpath)
	return path.Join(dir, file)
}
