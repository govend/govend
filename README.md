govend [![GoDoc](http://godoc.org/github.com/govend/govend?status.png)](http://godoc.org/github.com/govend/govend) [![Build Status](https://travis-ci.org/govend/govend.svg?branch=master)](https://travis-ci.org/govend/govend) [![Go Report Card](http://goreportcard.com/badge/govend/govend?)](http://goreportcard.com/report/govend/govend) [![Join the chat at https://gitter.im/govend/govend](https://badges.gitter.im/govend/govend.svg)](https://gitter.im/govend/govend?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
============================================================================================================================

`govend` leverages the `GO15VENDOREXPERIMENT` to vendor external package dependencies.

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
github.com/spf13/cobra
github.com/spf13/pflag
github.com/inconshreveable/mousetrap
github.com/cpuguy83/go-md2man
github.com/russross/blackfriday
github.com/shurcooL/sanitized_anchor_name
gopkg.in/yaml.v2
gopkg.in/check.v1
github.com/BurntSushi/toml
```

# Vendor Lock

The command `govend -v` only scans for external packages and downloads them to the `vendor` directory in your project. You may need more control over versioning your dependencies so that reliable reproducible builds are possible.

`govend` can save the path and commit revisions of each package dependency in a `vendor.yml` file in the root directory of your project. The format of the file can be `JSON` or `TOML` as well.

To lock in dependency versions run `govend -l` for `lock`.  An example `vendor.yml` file is shown below:

```yaml
vendors:
- path: github.com/BurntSushi/toml
  rev: 056c9bc7be7190eaa7715723883caffa5f8fa3e4
- path: github.com/cpuguy83/go-md2man
  rev: 71acacd42f85e5e82f70a55327789582a5200a90
- path: github.com/inconshreveable/mousetrap
  rev: 76626ae9c91c4f2a10f34cad8ce83ea42c93bb75
- path: github.com/kr/fs
  rev: 2788f0dbd16903de03cb8186e5c7d97b69ad387b
- path: github.com/russross/blackfriday
  rev: 300106c228d52c8941d4b3de6054a6062a86dda3
- path: github.com/shurcooL/sanitized_anchor_name
  rev: 10ef21a441db47d8b13ebcc5fd2310f636973c77
- path: github.com/spf13/cobra
  rev: 1c44ec8d3f1552cac48999f9306da23c4d8a288b
- path: github.com/spf13/pflag
  rev: 08b1a584251b5b62f458943640fc8ebd4d50aaa5
- path: gopkg.in/check.v1
  rev: 11d3bc7aa68e238947792f30573146a3231fc0f1
- path: gopkg.in/yaml.v2
  rev: 53feefa2559fb8dfa8d81baad31be332c97d6c77
```

# Report Summary
If you would like to get a report summary of the number of unique packages scanned, skipped and how many repositories were downloaded, run `govend -v -r`.

```bash
→ govend -v -r
github.com/BurntSushi/toml
gopkg.in/yaml.v2
gopkg.in/check.v1
github.com/kr/fs
github.com/spf13/cobra
github.com/spf13/pflag
github.com/cpuguy83/go-md2man
github.com/russross/blackfriday
github.com/shurcooL/sanitized_anchor_name
github.com/inconshreveable/mousetrap

packages scanned: 12
packages skipped: 0
repos downloaded: 10
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
`govend` works on Windows, but it may have some bugs.

# Contributing
Simply fork the code and send a pull request.
