package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const file = "vendor"

var extensions = []string{".json", ".yml", ".yaml", ".toml"}

// Manifest describes the vendors manifest file used save repository
// dependencies and their versions. The file is written as JSON, YAML or TOML.
type Manifest struct {
	format  string
	Vendors []vendor `json:"vendors" yml:"vendors" toml:"vendors"`
}

// vendor describes a dependecy with its import path and revision hash.
type vendor struct {
	Path string `json:"path" yaml:"path"`
	Rev  string `json:"rev,omitempty" yaml:"rev,omitempty"`
}

// SetFormat takes a string format and sets it, if it is valid.
// If the format provided is not a valid or supported format an error is
// returned.
//
// Currently YAML, JSON, and TOML are supported formats.
func (m *Manifest) SetFormat(format string) error {

	if format == "" {
		return errors.New("cannot set empty format")
	}

	format = strings.ToLower(format)
	for _, ext := range extensions {
		if format == string([]rune(ext)[1:]) {
			if format == "yaml" {
				format = "yml"
			}
			m.format = format
			return nil
		}
	}

	return fmt.Errorf("format type '%s' not supported", m.format)
}

// Format returns the manifest's set format.
func (m *Manifest) Format() string {
	return m.format
}

// Append creates a vendor object from a path and revision and
// appends it to the Manifest.
func (m *Manifest) Append(path, rev string) {
	for _, vendor := range m.Vendors {
		if vendor.Path == path {
			vendor.Rev = rev
			return
		}
	}
	m.Vendors = append(m.Vendors, vendor{path, rev})
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

// inSync check if the manifest file's list of vendored directories are on disk.
func (m Manifest) inSync() ([]vendor, bool) {
	inSync := true
	orphans := []vendor{}
	for _, v := range m.Vendors {
		if _, err := os.Stat(filepath.Join("vendor", v.Path)); os.IsNotExist(err) {
			orphans = append(orphans, v)
			inSync = false
		}
	}
	return orphans, inSync
}

// Sync removes orphaned vendored packages from the Manifest.
func (m *Manifest) Sync() {
	vendors, ok := m.inSync()
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
