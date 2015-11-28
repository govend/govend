package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Load reads a vendor manifest file and returns a Manifest object.
// It also takes an optional format to default the manifest to when written.
func Load(format string) (*Manifest, error) {

	// if the format is not specified default to YAML
	if format == "" {
		format = "yml"
	}

	// create a new manifest and set the format
	m := &Manifest{}
	if err := m.SetFormat(format); err != nil {
		return nil, err
	}

	// check if a yml, json, or toml manifest file exists on disk
	var f string
	for _, ext := range extensions {
		if _, err := os.Stat(file + ext); err == nil {
			f = string([]rune(ext)[1:])
		}
	}

	// if no format has been specified, no valid manifest file is present
	if f == "" {
		return m, nil
	}

	// read the manifest file on disk
	bytes, err := ioutil.ReadFile(file + "." + m.format)
	if err != nil {
		return m, err
	}

	switch m.format {
	case "json":
		if err := json.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
	case "yml":
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
	case "toml":
		if err := toml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("file type %s not supported", format)
	}

	return m, nil
}
