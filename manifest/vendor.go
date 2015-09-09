package manifest

// Vendor describes a vendored dependecy.
// vendor contains the import path and revision hash.
type Vendor struct {
	Path string `json:"path" yaml:"path"`
	Rev  string `json:"rev,omitempty" yaml:"rev,omitempty"`
}

func NewVendor(path, rev string) Vendor {
	return Vendor{path, rev}
}
