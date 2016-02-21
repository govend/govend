package packages

import (
	"os"
	"strings"

	"github.com/govend/govend/packages/filters"
	"github.com/kr/fs"
)

// Scan takes a directory and scans it for import dependencies.
func Scan(path string, pkg, testfiles, all bool) ([]string, error) {

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

		// skip directories named "vendor" or "testdata"
		if finfo.IsDir() {
			if finfo.Name() == "vendor" || finfo.Name() == "testdata" || pkg {
				w.SkipDir()
				continue
			}
		}

		// if testfiles is false then skip all go tests deps
		if !testfiles && strings.HasSuffix(finfo.Name(), "_test.go") {
			continue
		}

		// only parse .go files
		fpath := w.Path()
		if strings.HasSuffix(fpath, ".go") {
			imports, err := Parse(w.Path())
			if err != nil {

				// if the error is because of a bad file, skip the file
				if strings.Contains(err.Error(), eofError) {
					continue
				}
			}
			pkgs = append(pkgs, imports...)
		}
	}

	// filter packages
	if !all {
		pkgs = filters.Exceptions(pkgs)
		pkgs = filters.Standard(pkgs)
		pkgs = filters.Local(pkgs)
	}
	pkgs = filters.Duplicates(pkgs)

	return pkgs, nil
}
