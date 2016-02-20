![govend](art/govend.png)

# govend [![GoDoc](http://godoc.org/github.com/govend/govend?status.png)](http://godoc.org/github.com/govend/govend) [![Build Status](https://travis-ci.org/govend/govend.svg?branch=master)](https://travis-ci.org/govend/govend) [![Go Report Card](http://goreportcard.com/badge/govend/govend?1)](http://goreportcard.com/report/govend/govend) [![Join the chat at https://gitter.im/govend/govend](https://badges.gitter.im/govend/govend.svg)](https://gitter.im/govend/govend?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) ![](https://img.shields.io/badge/windows-ready-green.svg)

`govend` is a simple tool written in Golang to vendor Go packages as external or third party dependencies.

**it does:**
* try to be compatible with any project
* take a both note and code from `go get`
* vendor the nested dependency tree to the `nth` degree
* utilize the `GO15VENDOREXPERIMENT` and `vendor` directory as specified in golang version 1.5

**it does not:**
* wrap the `go` command
* try to create a new project
* generate temporary directories or files
* alter any go environment variables, including `$GOPATH`

# Supported Go Versions

* Go 1.4 or less - Go does not support vendoring
* Go 1.5 - vendor via `GO15VENDOREXPERIMENT=1`
* Go 1.6 - vendor unless `GO15VENDOREXPERIMENT=0`
* Go 1.7+ - vendor always despite the value of `GO15VENDOREXPERIMENT`

For further explanation please read https://golang.org/doc/go1.6#go_command.

# Install

```bash
$ go get -u github.com/govend/govend
```

# Vendor

To vendor external package dependencies run `govend` while inside the root directory of the project.  If you would like to see more verbose output run `govend -v`.

```bash
$ cd project/root

$ govend -v
github.com/kr/fs
github.com/BurntSushi/toml
github.com/spf13/cobra
github.com/inconshreveable/mousetrap
github.com/spf13/pflag
gopkg.in/yaml.v2
gopkg.in/check.v1
```

# Vendor Lock

The command `govend -v` only scans for external packages and downloads them to the `vendor` directory in your project. You may need more control over versioning your dependencies so that reliable reproducible builds are possible.

`govend` can save the path and commit revisions of each package dependency in a `vendor.yml` file in the root directory of your project. The format of the file can be `JSON` or `TOML` as well.

To lock in dependency versions run `govend -l` for `lock`.  An example `vendor.yml` file is shown below:

```yaml
vendors:
- path: github.com/BurntSushi/toml
  rev: f772cd89eb0b33743387f826d1918df67f99cc7a
- path: github.com/inconshreveable/mousetrap
  rev: 76626ae9c91c4f2a10f34cad8ce83ea42c93bb75
- path: github.com/kr/fs
  rev: 2788f0dbd16903de03cb8186e5c7d97b69ad387b
- path: github.com/spf13/cobra
  rev: 65a708cee0a4424f4e353d031ce440643e312f92
- path: github.com/spf13/pflag
  rev: 7f60f83a2c81bc3c3c0d5297f61ddfa68da9d3b7
- path: gopkg.in/check.v1
  rev: 4f90aeace3a26ad7021961c297b22c42160c7b25
- path: gopkg.in/yaml.v2
  rev: f7716cbe52baa25d2e9b0d0da546fcf909fc16b4
```

# Report Summary
If you would like to get a report summary of the number of unique packages scanned, skipped and how many repositories were downloaded, run `govend -v -r`.

```bash
→ govend -v -r
github.com/kr/fs
github.com/BurntSushi/toml
github.com/spf13/cobra
github.com/inconshreveable/mousetrap
github.com/spf13/pflag
gopkg.in/yaml.v2
gopkg.in/check.v1

packages scanned: 7
packages skipped: 0
repos downloaded: 7
```

# Scan
You may want to scan your code to determine how many third party dependencies are
in your project. To do so run `govend -s <path/to/dir>`. You can also specify different output formats.

**TXT**
```bash
$ govend -s packages
github.com/kr/fs
gopkg.in/yaml.v2
```

**JSON**
```bash
$ govend -s -f json packages
[
  "github.com/kr/fs",
  "gopkg.in/yaml.v2"
]
```

**YAML**
```bash
$ govend -s -f yaml packages
- github.com/kr/fs
- gopkg.in/yaml.v2
```
**XML**
```bash
$ govend -s -f xml packages
<string>gopkg.in/yaml.v2</string>
<string>github.com/kr/fs</string>
```

# More Flags

You can run `govend -h` to find more flags and options.

```bash
$ govend -h  
Govend downloads and vendors the packages named by the import
paths, along with their dependencies.

Usage:
  govend [flags]

Flags:
  -a, --all[=false]: The -a flag works with the -s flag to show all packages, not just
	external packages.

  -x, --commands[=false]: The -x flag prints commands as they are executed for vendoring
	such as 'git init'.

  -f, --format="YAML": The -f flag works with the -m flag and -s flag to define the
	format when writing files to disk.  By default, the file format is YAML but
	also supports JSON and TOML formats.

  -l, --lock[=false]: The -l flag writes a manifest vendor file on disk to lock in the
	versions of vendored dependencies.  This only needs to be done once.

  -r, --results[=false]: The -r flag works with the -v flag to print a summary of the
	number of packages scanned, packages skipped, and repositories downloaded.

  -s, --scan[=false]: The -s flag scans the current or provided directory for external
	packages.

  -u, --update[=false]: The -u flag uses the network to update the named packages and
	their dependencies.  By default, the network is used to check out missing
	packages but does not use it to look for updates to existing packages.

  -v, --verbose[=false]: The -v flag prints the names of packages as they are vendored.

      --version[=false]: The --version flag prints the current version.
```

# Windows Support
`govend` works on Windows, but please report any bugs.

# Contributing
Simply fork the code and send a pull request.
