// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package filters

// ExceptionList is a list of exception packages to skip.
var ExceptionList = []string{
	"appengine",
	"appengine/*",
	"appengine_internal",
	"appengine_internal/*",
}
