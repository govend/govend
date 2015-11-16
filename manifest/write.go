package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Write writes the vendors to the manifest file on disk.
func (m *Manifest) Write() error {

	var b []byte
	var err error
	var fpath string

	// sort the manfiest vendors to fixate ordering
	sort.Sort(m)

	switch format {
	case "json":
		fpath = file + "." + format
		b, err = json.Marshal(m)
		if err != nil {
			return err
		}
	case "yml", "yaml":
		format = "yml"
		fpath = file + "." + format
		b, err = yaml.Marshal(m)
		if err != nil {
			return err
		}
	case "toml":
		fpath = file + "." + format
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(m); err != nil {
			log.Fatal(err)
		}
		b = buf.Bytes()
	default:
		return fmt.Errorf("vendor manifest file format type '%s' is not supported", format)
	}

	// cleanup any previous file formats
	if err := os.Remove(file + ".yml"); err != nil {
		return err
	}
	if err := os.Remove(file + ".json"); err != nil {
		return err
	}
	if err := os.Remove(file + ".toml"); err != nil {
		return err
	}

	// write manifest file bytes to disk
	if err := ioutil.WriteFile(fpath, b, 0644); err != nil {
		return err
	}

	return nil
}
