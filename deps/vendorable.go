// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/govend/govend/deps/semver"
)

func Vendorable() error {
	go15, _ := semver.New("1.5.0")
	go16, _ := semver.New("1.6.0")
	go17, _ := semver.New("1.7.0")

	version, err := semver.New(strings.TrimPrefix(runtime.Version(), "go"))
	if err != nil {
		return err
	}

	if version.LessThan(go15) {
		return errors.New("vendoring requires Go versions 1.5+")
	}

	if version.GreaterThanEqual(go15) && version.LessThan(go16) {
		if os.Getenv("GO15VENDOREXPERIMENT") != "1" {
			return errors.New("Go 1.5.x requires 'GO15VENDOREXPERIMENT=1'")
		}
	}

	if version.GreaterThanEqual(go16) && version.LessThan(go17) {
		if os.Getenv("GO15VENDOREXPERIMENT") == "0" {
			return errors.New("Go 1.6.x cannot vendor with 'GO15VENDOREXPERIMENT=0'")
		}
	}

	return nil
}
