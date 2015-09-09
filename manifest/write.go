package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Write writes the vendors to the manifest file on disk.
func Write(file string, vendors *[]Vendor) error {

	var bytes []byte
	var err error

	// marshal by format type
	switch format {

	// marshal JSON
	case "json":
		bytes, err = json.Marshal(&vendors)
		if err != nil {
			return err
		}

		// marshal YML
	case "yml", "yaml", "":
		bytes, err = yaml.Marshal(&vendors)
		if err != nil {
			return err
		}

	// error on unsupported extension type.
	default:
		return fmt.Errorf("vendor manifest file format type '%s' is not supported", format)
	}

	// write file
	if err := ioutil.WriteFile(file, bytes, 0777); err != nil {
		return err
	}

	return nil
}
