// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements access to gc-generated export data.

package main

import (
	"github.com/gophersaurus/govend/internal/_vendor/golang.org/x/tools/go/gcimporter"
)

func init() {
	register("gc", gcimporter.Import)
}
