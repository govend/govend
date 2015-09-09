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

// Scan scans a directory to collect external package imports.
func Scan(dir, vendorDir string, skipVendor bool) ([]string, error) {

	var pkglist []string

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

		// check the file is a .go file
		if strings.HasSuffix(w.Path(), ".go") {

			fset := token.NewFileSet()

			// parse only the import declarations in the .go file
			f, err := parser.ParseFile(fset, w.Path(), nil, parser.ImportsOnly)
			if err != nil {

				msg := "expected 'package', found 'EOF'"
				e := err.Error()
				if len(e) >= len(msg) {
					if e[len(e)-len(msg):] == msg {
						continue
					}
				}
				return nil, err
			}

			// iterate over import paths
			for _, i := range f.Imports {

				// unquote the import path value
				importpath, err := strconv.Unquote(i.Path.Value)
				if err != nil {
					return nil, err
				}

				// iterate through the known external packages
				for _, pkg := range pkglist {
					if importpath == pkg {
						goto SKIP
					}
				}

				pkglist = append(pkglist, importpath)

			SKIP: // skips the appending of packages that are already present
			}
		}
	}

	return pkglist, nil
}
