package go15experiment

import (
	"os"
	"os/exec"
)

// Version returns true of the Go version is correct for the GO15VENDOREXPERIMENT.
func Version() bool {

	// 6l was removed in 1.5, when vendoring was introduced.
	cmd := exec.Command("go", "tool", "6l")
	if _, err := cmd.CombinedOutput(); err == nil {
		return false
	}
	return true
}

// On returns true of the GO15VENDOREXPERIMENT enviorment variable is on.
func On() bool {
	return os.Getenv("GO15VENDOREXPERIMENT") == "1"
}
