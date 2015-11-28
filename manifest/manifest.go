package manifest

import (
	"os"
	"path/filepath"
)

const file = "vendor"

// Manifest describes the vendors manifest file used save repository
// dependencies and their versions. The file is written as JSON, YAML or TOML.
type Manifest struct {
	Vendors []Vendor `json:"vendors" yml:"vendors" toml:"vendors"`
}

// Append takes a Vendor and appends it to the Manifest object.
func (m *Manifest) Append(v Vendor) {
	m.Vendors = append(m.Vendors, v)
}

// Remove takes a package import string, and removes it from the manifest file.
func (m *Manifest) Remove(pkg string) {
	for key, vendor := range m.Vendors {
		if vendor.Path == pkg {
			m.Vendors = append(m.Vendors[:key], m.Vendors[key+1:]...)
		}
	}
}

// Contains returns true of the package import string is contained in the
// Manifest object
func (m Manifest) Contains(pkg string) bool {
	for _, vendor := range m.Vendors {
		if vendor.Path == pkg {
			return true
		}
	}
	return false
}

// InSync check if the manifest file's list of vendored directories are on disk.
func (m Manifest) InSync() ([]Vendor, bool) {
	inSync := true
	var orphans []Vendor
	for _, vendor := range m.Vendors {
		if _, err := os.Stat(filepath.Join("vendor", vendor.Path)); os.IsNotExist(err) {
			orphans = append(orphans, vendor)
			inSync = false
		}
	}
	return orphans, inSync
}

// Sync acts like InSync, but removes orphaned vendored packages from the
// Manifest object.
func (m *Manifest) Sync() {
	vendors, ok := m.InSync()
	if !ok {
		for _, vendor := range vendors {
			m.Remove(vendor.Path)
		}
	}
}

// Len allows Manifest to satisfy the sort.Interface.
func (m *Manifest) Len() int {
	return len(m.Vendors)
}

// Swap allows Manifest to satisfy the sort.Interface.
func (m *Manifest) Swap(i, j int) {
	m.Vendors[i], m.Vendors[j] = m.Vendors[j], m.Vendors[i]
}

// Less allows Manifest to satisfy the sort.Interface.
func (m *Manifest) Less(i, j int) bool {
	return m.Vendors[i].Path < m.Vendors[j].Path
}

// Vendor describes an external vendored dependecy with its import path and
// revision hash.
type Vendor struct {
	Path string `json:"path" yaml:"path"`
	Rev  string `json:"rev,omitempty" yaml:"rev,omitempty"`
}

func NewVendor(path, rev string) Vendor {
	return Vendor{path, rev}
}
