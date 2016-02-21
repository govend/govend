//go:generate go run generate.go

package filters

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
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

	prefix := projectImportPath()
	l := len(prefix)

	var list []string
	for _, pkg := range pkgs {

		// check the item length is geater than or equal to the prefix length
		// this ensures no out of bounds memory errors will occur
		if len(pkg) >= l {
			if prefix == pkg[:l] {
				continue
			}
		}

		list = append(list, pkg)
	}

	return list
}

// projectImportPath returns the import path of the current project directory.
// It does so via $GOPATH.
func projectImportPath() string {

	// determine the current absolute file path
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// check the env $GOPATH is valid
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		log.Fatal(errors.New("$GOPATH not set"))
	}

	// leverage the $GOPATH to strip out everything but the base git URL
	return path[len(gopath+"/src/"):]
}
