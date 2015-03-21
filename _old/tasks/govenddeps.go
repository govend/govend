package tasks

import (
	"os/exec"

	"github.com/jackspirou/govend/app/helpers"
)

func govendDeps() {

	// deps for govend.
	err := exec.Command("go", "get", "-d", "github.com/jackspirou/rmimports").Run()
	helpers.Check(err)

	// deps for govend.
	err = exec.Command("go", "get", "github.com/jackspirou/gormimports").Run()
	helpers.Check(err)

}
