// +build ignore

// Command mkindex creates the file "stdpkgs.go" containing an index of the Go
// standard library.

// This file has been modified from is source which can be found here:
// https://github.com/golang/tools/blob/master/imports/mkindex.go

// To execute this file and generate stdpkgs.go, I simply "go run mkstdpkgs.go".
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

var (
	pkgIndex = make(map[string][]pkg)
)

func main() {
	// Don't use GOPATH.
	ctx := build.Default
	ctx.GOPATH = ""

	// Populate pkgIndex global from GOROOT.
	for _, path := range ctx.SrcDirs() {
		f, err := os.Open(path)
		if err != nil {
			log.Print(err)
			continue
		}
		children, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Print(err)
			continue
		}
		for _, child := range children {
			if child.IsDir() {
				loadPkg(path, child.Name())
			}
		}
	}

	// Construct source file.
	var buf bytes.Buffer
	fmt.Fprint(&buf, pkgIndexHead)
	fmt.Fprintf(&buf, "var stdpkgs = %#v\n", pkgIndex)
	src := buf.Bytes()

	// Replace main.pkg type name with pkg.
	src = bytes.Replace(src, []byte("main.pkg"), []byte("pkg"), -1)
	// Replace actual GOROOT with "/go".
	src = bytes.Replace(src, []byte(ctx.GOROOT), []byte("/go"), -1)
	// Add some line wrapping.
	src = bytes.Replace(src, []byte("}, "), []byte("},\n"), -1)
	src = bytes.Replace(src, []byte("true, "), []byte("true,\n"), -1)

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

const pkgIndexHead = `package main
`

var fset = token.NewFileSet()

func loadPkg(root, importpath string) {
	shortName := path.Base(importpath)
	if shortName == "testdata" {
		return
	}

	dir := filepath.Join(root, importpath)
	pkgIndex[shortName] = append(pkgIndex[shortName], pkg{
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
			loadPkg(root, filepath.Join(importpath, name))
		}
	}
}

type pkg struct {
	importpath string // full pkg import path, e.g. "net/http"
	dir        string // absolute file path to pkg directory e.g. "/usr/lib/go/src/fmt"
}
