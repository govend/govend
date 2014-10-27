package tasks

import "os/exec"

func goRmImports() {
	// Run gormimports
	exec.Command("bash", "-c", "gormimports -w ./*").Run()
}
