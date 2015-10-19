package packages

//go:generate go run generate_std_pkgs.go

import (
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"

	"github.com/kr/fs"
)

const (
	eofError        = "expected 'package', found 'EOF'"
	importPathError = "invalid import path:"
)

var (
	rootExceptions = []string{"appengine/"}
	scannedFiles   = map[string][]string{}
)

// Scan scans a directory to collect external package imports.
func Scan(dir, vendorDir string, skipVendor, verbose bool) ([]string, map[string]string, error) {

	pkglist := []string{}
	badpkgs := map[string]string{}

	w := fs.Walk(dir)
	for w.Step() {

		fstat := w.Stat()

		if fstat.IsDir() {

			if fstat.Name() == "testdata" {
				w.SkipDir()
				continue
			}

			// check if that directory is "_vendor"
			if fstat.Name() == vendorDir && skipVendor {
				w.SkipDir()
				continue
			}
			continue
		}

		// check for errors
		if w.Err() != nil {
			log.Println("govend scan:", w.Err())
			continue
		}

		fpath := w.Path()

		// check the file is a .go file
		if strings.HasSuffix(fpath, ".go") {

			var imports []string

			if pkgs, ok := scannedFiles[fpath]; ok {
				imports = pkgs
			} else {

				fset := token.NewFileSet()

				// parse only the import declarations in the .go file
				f, err := parser.ParseFile(fset, w.Path(), nil, parser.ImportsOnly)
				if err != nil {

					e := err.Error()

					if strings.Contains(e, eofError) {
						continue
					}

					if strings.Contains(e, importPathError) {
						for _, i := range f.Imports {
							badpkgs[i.Path.Value] = "invalid import path"
						}
						continue
					}

					return nil, badpkgs, err
				}

				// unquote the import path value
				for _, i := range f.Imports {
					importpath, err := strconv.Unquote(i.Path.Value)
					if err != nil {
						return nil, badpkgs, err
					}
					imports = append(imports, importpath)
				}
			}

			// iterate over import paths
			for _, importpath := range imports {

				// iterate through the known external packages
				for _, pkg := range pkglist {
					if importpath == pkg {
						goto SKIP
					}
				}

				for _, root := range rootExceptions {
					if len(importpath) >= len(root) {
						if importpath[:len(root)] == root {
							goto SKIP
						}
					}
				}

				pkglist = append(pkglist, importpath)
			SKIP: // skips the appending of packages that are already present
			}
		}
	}

	return pkglist, badpkgs, nil
}
