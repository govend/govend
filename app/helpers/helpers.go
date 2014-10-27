package helpers

import (
	"go/scanner"
	"log"
	"os"
	"strings"
)

// Check takes an error, prints it, and then exits with 1.
func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Report(err error) {
	scanner.PrintError(os.Stderr, err)
}

func IsGoFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}
