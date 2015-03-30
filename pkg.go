package main

// pkg represents a package in "stdpkgs.go".
type pkg struct {
	path string // full pkg import path, e.g. "net/http"
	dir  string // absolute file path to pkg directory e.g. "/usr/lib/go/src/fmt"
}
