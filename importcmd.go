package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jackspirou/pimports"
)

var (
	options = &pimports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}

	list   bool
	write  bool
	doDiff bool

	priority string
)

func importcmd(d, e, l, w bool, p string, args []string) error {

	list = l
	write = w
	doDiff = d

	options.AllErrors = e

	priority = p

	projectpath, err := importpath(".")
	if err != nil {
		return err
	}

	projectpath = projectpath + "/" + vendorDir

	if len(priority) < 1 {
		priority = projectpath
	} else {
		priority = projectpath + "," + priority
	}

	if len(priority) > 0 {
		ptemp := strings.Split(priority, ",")
		options.Priority = ptemp
	}

	if options.TabWidth < 0 {
		return fmt.Errorf("negative tabwidth %d\n", options.TabWidth)
	}

	if len(args) < 1 {
		args = []string{"."}
	}

	for _, v := range args {
		switch dir, err := os.Stat(v); {
		case err != nil:
			return err
		case dir.IsDir():
			importWalkDir(v)
		default:
			if err := importProcessFile(v, nil, os.Stdout, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func importProcessFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	opt := options
	if stdin {
		nopt := *options
		nopt.Fragment = true
		opt = &nopt
	}

	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := pimports.Process(filename, src, opt)
	if err != nil {
		return err
	}

	if !bytes.Equal(src, res) {
		// formatting has changed
		if list {
			fmt.Fprintln(out, filename)
		}
		if write {
			err = ioutil.WriteFile(filename, res, 0)
			if err != nil {
				return err
			}
		}
		if doDiff {
			data, err := importDiff(src, res)
			if err != nil {
				return fmt.Errorf("computing diff: %s", err)
			}
			fmt.Printf("diff %s gofmt/%s\n", filename, filename)
			out.Write(data)
		}
	}

	if !list && !write && !doDiff {
		_, err = out.Write(res)
	}

	return err
}

func importVisitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = importProcessFile(path, nil, os.Stdout, false)
	}
	if err != nil {
		return err
	}
	return nil
}

func importWalkDir(path string) {
	filepath.Walk(path, importVisitFile)
}

func importDiff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "gofmt")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "gofmt")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}
