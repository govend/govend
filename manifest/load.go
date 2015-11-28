package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

var format string

// Load reads a vendor manifest file and returns a Manifest object.
func Load() (*Manifest, error) {

	// check if a yml, json, or toml manifest file exists
	if _, err := os.Stat(file + ".yml"); err == nil {
		format = "yml"
	}
	if _, err := os.Stat(file + ".json"); err == nil {
		format = "json"
	}
	if _, err := os.Stat(file + ".xml"); err == nil {
		format = "xml"
	}

	// define an empty Manifest file to use as state
	m := &Manifest{}

	// if no format has been specified, no valid manifest file is present
	// default to yaml
	if format == "" {
		format = "yml"
		return m, nil
	}

	// define the file path
	fpath := file + "." + format

	// read the manifest file on disk
	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return m, err
	}

	switch format {
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
		return nil, fmt.Errorf("vendor manifest file type %s not supported", format)
	}

	return m, nil
}
