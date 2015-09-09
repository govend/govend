package repo

import "fmt"

// Ping checks the host server to determine the package Repo.
func Ping(pkg string) (*Repo, error) {

	// determine import path and repository type by pinging the server host
	repo, err := ImportPath(pkg, false)
	if err != nil {
		e := err.Error()
		msg := "no go-import meta tags"
		if e[len(e)-len(msg):] == msg {
			return nil, fmt.Errorf("network ping (potential proxy issue): %s ", e)
		}
		return nil, fmt.Errorf("network ping: %s", err)
	}

	return repo, nil
}
