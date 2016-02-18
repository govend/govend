package semver

import (
	"strconv"
	"strings"
)

type SemVer struct {
	Major int
	Minor int
	Patch int
}

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
