package packages

import (
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/gophersaurus/govend/strutil"
	"github.com/kr/fs"
)

// ScanProject
func ScanProject(dir string) ([]string, error) {

	imports := map[string]bool{}

	w := fs.Walk(dir)
	for w.Step() {

		fstat := w.Stat()

		if fstat.IsDir() {

			// check if that directory is "_vendor"
			n := fstat.Name()
			if n == "vendor" || n == "testdata" || []rune(n)[0] == '_' {
				w.SkipDir()
				continue
			}
		}

		// check for errors
		if w.Err() != nil {
			return nil, w.Err()
		}

		fpath := w.Path()

		// check the file is a .go file
		if strings.HasSuffix(fpath, ".go") {

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
						if Valid(i.Path.Value) {
							path, err := strconv.Unquote(i.Path.Value)
							if err != nil {
								return nil, err
							}
							imports[path] = true
							continue
						}
					}
					continue
				}

				return nil, err
			}

			// unquote the import path value
			for _, i := range f.Imports {
				path, err := strconv.Unquote(i.Path.Value)
				if err != nil {
					return nil, err
				}
				imports[path] = true
			}
		}
	}

	for path := range imports {
		for _, exception := range Exceptions {
			if Match(path, exception) {
				imports[path] = false
			}
		}
	}

	paths := []string{}
	for path, ok := range imports {
		if ok {
			paths = append(paths, path)
		}
	}

	paths = FilterStdPkgs(paths)

	projectpath, err := ImportPath(".")
	if err != nil {
		return nil, err
	}

	// filter out packages internal to the project
	paths = strutil.RemovePrefixInStringSlice(projectpath, paths)

	return paths, nil
}
