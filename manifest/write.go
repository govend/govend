package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

// Write writes the vendors to the manifest file on disk.
func Write(file string, vendors *[]Vendor) error {

	var bytes []byte
	var err error

	// sort vendors to fixate ordering in manifest file
	sort.Sort(sorter(*vendors))

	// marshal by format type
	switch format {

	case "json":
		bytes, err = json.Marshal(&vendors)
		if err != nil {
			return err
		}
	case "yml", "yaml", "":
		bytes, err = yaml.Marshal(&vendors)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("vendor manifest file format type '%s' is not supported", format)
	}

	if err := ioutil.WriteFile(file, bytes, 0644); err != nil {
		return err
	}

	return nil
}
