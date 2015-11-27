package packages

//go:generate go run generate_std_pkgs.go

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/kr/fs"
)

const (
	eofError        = "expected 'package', found 'EOF'"
	importPathError = "invalid import path:"
)

func Scan(dir string) ([]string, error) {

	fileInfo, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("'%s' is a file, directories contain packages", fileInfo.Name())
	}

	imports := map[string]bool{}

	w := fs.Walk(dir)
	w.Step()

	for w.Step() {

		fstat := w.Stat()

		if fstat.IsDir() {
			w.SkipDir()
			continue
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

	for path, _ := range imports {
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

	return paths, nil
}
