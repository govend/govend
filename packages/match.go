package packages

import "strings"

// Match determins if two package paths match.
func Match(path1, path2 string) bool {

	p1 := strings.Split(path1, "/")
	p2 := strings.Split(path2, "/")

	for i, v := range p1 {
		if i < len(p2) {
			if v != p2[i] {
				if v == "*" || p2[i] == "*" {
					return true
				}
				return false
			}
		}
	}

	if len(p1) != len(p2) {
		return false
	}

	return true
}
