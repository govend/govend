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
