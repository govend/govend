// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package filters

import "strings"

// Match determines if two package paths match.
func Match(path1, path2 string) bool {

	p1 := strings.Split(path1, "/")
	p2 := strings.Split(path2, "/")

	for i, v := range p1 {
		if i < len(p2) {
			if v != p2[i] {
				return v == "*" || p2[i] == "*"
			}
		}
	}

	return len(p1) == len(p2)
}
