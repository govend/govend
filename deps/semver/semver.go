// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

// Package semver provides a way to parse and compare semantic versions.
package semver

import (
	"strconv"
	"strings"
)

// SemVer describes a semantic version number.
type SemVer struct {
	Major int
	Minor int
	Patch int
}

// New takes a semantic version number as a string and returns a SemVer object.
func New(verison string) (SemVer, error) {
	sem := strings.Split(verison, ".")
	switch len(sem) {
	case 3:
		major, err := strconv.Atoi(sem[0])
		if err != nil {
			return SemVer{}, err
		}
		minor, err := strconv.Atoi(sem[1])
		if err != nil {
			return SemVer{}, err
		}
		patch, err := strconv.Atoi(sem[2])
		if err != nil {
			return SemVer{}, err
		}
		return SemVer{Major: major, Minor: minor, Patch: patch}, nil
	case 2:
		major, err := strconv.Atoi(sem[0])
		if err != nil {
			return SemVer{}, err
		}
		minor, err := strconv.Atoi(sem[1])
		if err != nil {
			return SemVer{}, err
		}
		return SemVer{Major: major, Minor: minor, Patch: 0}, nil
	default:
		return SemVer{}, nil
	}
}

// GreaterThan takes a SemVer object and reports true if it is greater than
// the SemVer object.
func (s SemVer) GreaterThan(sem SemVer) bool {
	if s.Major > sem.Major {
		return true
	}
	if s.Major == sem.Major {
		if s.Minor > sem.Minor {
			return true
		}
		if s.Minor == sem.Minor {
			if s.Patch > sem.Patch {
				return true
			}
		}
	}
	return false
}

// GreaterThanEqual takes a SemVer object and reports true if it is greater
// than or equal to the SemVer object.
func (s SemVer) GreaterThanEqual(sem SemVer) bool {
	if s.Major > sem.Major {
		return true
	}
	if s.Major == sem.Major {
		if s.Minor > sem.Minor {
			return true
		}
		if s.Minor == sem.Minor {
			if s.Patch >= sem.Patch {
				return true
			}
		}
	}
	return false
}

// LessThan takes a SemVer object and reports true if it is less than
// the SemVer object.
func (s SemVer) LessThan(sem SemVer) bool {
	if s.Major < sem.Major {
		return true
	}
	if s.Major == sem.Major {
		if s.Minor < sem.Minor {
			return true
		}
		if s.Minor == sem.Minor {
			if s.Patch < sem.Patch {
				return true
			}
		}
	}
	return false
}
