// +build ignore

// This file generates stdpkgs.go, which contains the standard library packages.

// This file has been modified from its original source:
// https://github.com/golang/tools/blob/master/imports/mkindex.go

package main

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

// pkgs is a map of the standard golang packages.
var pkgs = make(map[string][]pkg)

// fset is a set of file tokens.
var fset = token.NewFileSet()

// pkg represents a packages.
type pkg struct {
	importpath string // full pkg import path, e.g. "net/http"
	dir        string // absolute file path to pkg directory e.g. "/usr/lib/go/src/fmt"
}

func main() {

	// start with the default context.
	ctx := build.Default

	// remove the GOPATH, we only want to search packages in the GOROOT.
	ctx.GOPATH = ""

	// iterate through the list of package source root directories.
	for _, path := range ctx.SrcDirs() {

		// open the file
		f, err := os.Open(path)
		if err != nil {
			log.Print(err)
			continue
		}

		// grab all the child names from the directory in a single slice.
		children, err := f.Readdir(-1)

		// close the file and check for errors.
		f.Close()
		if err != nil {
			log.Print(err)
			continue
		}

		// iterate through each child name.
		for _, child := range children {

			// check the child name is a directory.
			if child.IsDir() {

				// load the package path and name.
				load(path, child.Name())
			}
		}
	}

	// start with a byte buffer.
	var buf bytes.Buffer

	// write preliminary data to the byte buffer.
	buf.WriteString("package main\n")
	fmt.Fprintf(&buf, "var stdpkgs = %#v\n", pkgs)

	// transfer buffer bytes to final source.
	src := buf.Bytes()

	// replace main.pkg type name with pkg.
	src = bytes.Replace(src, []byte("main.pkg"), []byte("pkg"), -1)

	// replace actual GOROOT with "/go".
	src = bytes.Replace(src, []byte(ctx.GOROOT), []byte("/go"), -1)

	// add line wrapping.
	src = bytes.Replace(src, []byte("[]pkg{"), []byte("[]pkg{\n"), -1)
	src = bytes.Replace(src, []byte(", pkg"), []byte(",\npkg"), -1)
	src = bytes.Replace(src, []byte("}}, "), []byte("},\n},\n"), -1)
	src = bytes.Replace(src, []byte("true, "), []byte("true,\n"), -1)
	src = bytes.Replace(src, []byte("}}}"), []byte("},\n},\n}"), -1)

	var err error
	src, err = format.Source(src)
	if err != nil {
		log.Fatal(err)
	}

	// Write out source file.
	err = ioutil.WriteFile("stdpkgs.go", src, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// load takes a path root and import path
func load(root, importpath string) {
	shortName := path.Base(importpath)
	if shortName == "testdata" {
		return
	}

	dir := filepath.Join(root, importpath)
	pkgs[shortName] = append(pkgs[shortName], pkg{
		importpath: importpath,
		dir:        dir,
	})

	pkgDir, err := os.Open(dir)
	if err != nil {
		return
	}
	children, err := pkgDir.Readdir(-1)
	pkgDir.Close()
	if err != nil {
		return
	}
	for _, child := range children {
		name := child.Name()
		if name == "" {
			continue
		}
		if c := name[0]; c == '.' || ('0' <= c && c <= '9') {
			continue
		}
		if child.IsDir() {
			load(root, filepath.Join(importpath, name))
		}
	}
}
