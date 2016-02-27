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

	// sort the manfiest vendors to fixate ordering
	sort.Sort(m)

	// set a slice of bytes aside for the manifest formatted contents
	b := []byte{}

	var err error
	switch m.fmt {
	case "json":
		b, err = json.Marshal(m)
		if err != nil {
			return err
		}
	case "yml":
		b, err = yaml.Marshal(m)
		if err != nil {
			return err
		}
	case "toml":
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(m); err != nil {
			log.Fatal(err)
		}
		b = buf.Bytes()
	default:
		return fmt.Errorf("format type %q not supported", m.fmt)
	}

	// cleanup any previous file formats
	for _, ext := range extensions {
		os.Remove(file + ext)
	}

	// write manifest file bytes to disk
	if err := ioutil.WriteFile(file+"."+m.fmt, b, 0644); err != nil {
		return err
	}

	return nil
}
