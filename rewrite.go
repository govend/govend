package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/jackspirou/govend/internal/_vendor/github.com/kr/fs"
)

// rewriteImports takes a directory path and a map of key value strings and
// updates the import path with those replacements.
func rewrite(dir string, replace map[string]string) error {

	// create a new walk.
	w := fs.Walk(dir)

	// start the walk down the directory tree.
	for w.Step() {

		// check for errors.
		if w.Err() != nil {
			log.Println("gf import rewrite:", w.Err())
			continue
		}

		// check the file is a .go file.
		if !w.Stat().IsDir() && strings.HasSuffix(w.Path(), ".go") {

			// rewrite the file.
			if err := rewriteFile(w.Path(), replace); err != nil {
				return err
			}
		}
	}
	return nil
}

// match takes a file path and replacement map.
func match(path string, replace map[string]string) (string, bool) {
	for key, value := range replace {
		result := strings.Replace(path, key, value, 1)
		if path != result {
			return result, true
		}
	}
	return path, false
}

// rewriteFile rewrites import statments in the named file
// according to the rules supplied by the map of strings.
func rewriteFile(name string, replace map[string]string) error {

	// create an empty fileset.
	fset := token.NewFileSet()

	// parse the .go file.
	// we are parsing the entire file with comments, so we don't lose anything
	// if we need to write it back out.
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		e := err.Error()
		msg := "expected 'package', found 'EOF'"
		if e[len(e)-len(msg):] == msg {
			return nil
		}
		return err
	}

	// iterate through the import paths. if a change occurs update bool.
	change := false
	for _, i := range f.Imports {

		// unquote the import path value.
		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			return err
		}

		// match import path with the given replacement map
		if path, ok := match(path, replace); ok {
			i.Path.Value = strconv.Quote(path)
			change = true
		}
	}

	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// import \"") {

				// trim off extra comment stuff
				ctext := c.Text
				ctext = strings.TrimPrefix(ctext, "// import")
				ctext = strings.TrimSpace(ctext)

				// unquote the comment import path value
				ctext, err := strconv.Unquote(ctext)
				if err != nil {
					return err
				}

				// match the comment import path with the given replacement map
				if ctext, ok := match(ctext, replace); ok {
					c.Text = "// import " + strconv.Quote(ctext)
					change = true
				}
			}
		}
	}

	// if no change occured, then we don't need to write to disk, just return.
	if !change {
		return nil
	}

	// since the imports changed, resort them.
	ast.SortImports(fset, f)

	// create a temporary file, this easily avoids conflicts.
	temp := name + ".temp"
	w, err := os.Create(temp)
	if err != nil {
		return err
	}

	// write changes to .temp file, and include proper formatting.
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(w, fset, f)
	if err != nil {
		return err
	}

	// close the writer
	err = w.Close()
	if err != nil {
		return err
	}

	// rename the .temp to .go
	return os.Rename(temp, name)
}
