// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package imports

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// Format takes a slice of packages and formats it.
func Format(pkgs []string, format string) ([]byte, error) {

	if format == "YAML" {
		format = ""
	}
	format = strings.ToLower(format)

	// create an slice of bytes to print or write results
	b := []byte{}
	var err error

	// switch on format
	switch format {

	case "txt", "": // if text, use a byte.Buffer to format package paths
		var buff bytes.Buffer
		for _, pkg := range pkgs {
			buff.WriteString(pkg + "\n")
		}
		b = buff.Bytes()

	case "xml": // marshal to xml with indentation
		b, err = xml.MarshalIndent(pkgs, "", "  ")
		if err != nil {
			return nil, err
		}

	case "yaml", "yml": // marshal to yml with indentation
		b, err = yaml.Marshal(pkgs)
		if err != nil {
			return nil, err
		}

	case "json": // marshal to json with indentation
		b, err = json.MarshalIndent(pkgs, "", "  ")
		if err != nil {
			return nil, err
		}

	default: // error out on unsupported formats
		return nil, fmt.Errorf("unsupported format '%s' ", format)
	}

	return b, nil
}
