package tasks

import "os/exec"

func goImports() {
	// Run goimports
	exec.Command("bash", "-c", "goimports -w ./*").Run()
}
