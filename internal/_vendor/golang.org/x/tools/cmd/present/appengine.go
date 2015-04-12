// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package main

import (
	"mime"

	_ "github.com/gophersaurus/govend/internal/_vendor/golang.org/x/tools/playground"
	"github.com/gophersaurus/govend/internal/_vendor/golang.org/x/tools/present"
)

var basePath = "./present/"

func init() {
	initTemplates(basePath)
	playScript(basePath, "HTTPTransport")
	present.PlayEnabled = true

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}
