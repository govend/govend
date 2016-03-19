// Copyright 2016 govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package repos

// Ping checks the host server to determine the package Repo.
func Ping(pkg string) (*Repo, error) {

	// determine import path and repository type by pinging the server host
	repo, err := ImportPath(pkg, false)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
