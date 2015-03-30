package main

import "path"

// removestdpkgs
func removestdpkgs(pkgs []string) []string {

	// define an empty package list to fill
	var pkglist []string

	// iterate through the files import paths
	for _, pkg := range pkgs {

		// determine the name of the package
		name := path.Base(pkg)

		// skip CGO and any relative import paths
		if pkg == "C" || pkg[0] == '.' {
			continue
		}

		// if the package is part of the golang standard library, skip it
		if stds, ok := stdpkgs[name]; ok {
			for _, stdpkg := range stds {
				if pkg == stdpkg.path {
					goto SKIP
				}
			}
		}

		// if the import path doens't exists in pkgs, add it
		pkglist = append(pkglist, pkg)

	SKIP: // skips the appending of packages that are already present
	}

	return pkglist

}
