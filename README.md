govend [![GoDoc](http://godoc.org/github.com/gophersaurus/govend?status.png)](http://godoc.org/github.com/gophersaurus/govend) [![Build Status](https://travis-ci.org/gophersaurus/govend.svg?branch=master)](https://travis-ci.org/gophersaurus/govend) [![Go Report Card](http://goreportcard.com/badge/gophersaurus/govend?)](http://goreportcard.com/report/gophersaurus/govend)
============================================================================================================================

`govend` leverages the `GO15VENDOREXPERIMENT` to vendor dependencies.

**it does:**
* try to be compatible with any project
* vendor nested dependencies to the `nth` degree
* use the `vendor` directory as specified in golang version 1.5
* provide `goimports` functionality, while prioritizing vendored packages via the `govend imports` command

**it does not:**
* wrap the `go` command
* create a new project for you
* alter any environment variables, including `$GOPATH`

Install
=======

```bash
$ go get -u github.com/gophersaurus/govend
```

Vendor Dependencies
===================

```bash
$ cd project/root

$ govend -v
```

To vendor dependencies run `govend -v` while in the project root directory.

```bash
→ govend -v
5 packages scanned, 5 repositories found
 ↓ gopkg.in/yaml.v2 (latest)
 ↓ github.com/jackspirou/importsplus (latest)
 ↓ golang.org/x/tools (latest)
 ↓ github.com/spf13/cobra (latest)
 ↓ github.com/kr/fs (latest)

downloading recursive dependencies...

55 packages scanned, 10 repositories found
 ↓ github.com/inconshreveable/mousetrap (latest)
 ↓ golang.org/x/net (latest)
 ↓ github.com/spf13/pflag (latest)
 ↓ github.com/cpuguy83/go-md2man (latest)
 ↓ gopkg.in/check.v1 (latest)

downloading recursive dependencies...

74 packages scanned, 12 repositories found
 ↓ golang.org/x/text (latest)
 ↓ github.com/russross/blackfriday (latest)

downloading recursive dependencies...

95 packages scanned, 13 repositories found
 ↓ github.com/shurcooL/sanitized_anchor_name (latest)
```

### `vendors.yml`

The `vendors.yml` file contains an array of import paths and commit revisions.

```yaml
- path: gopkg.in/yaml.v2
  rev: 7ad95dd0798a40da1ccdff6dff35fd177b5edf40
- path: github.com/jackspirou/importsplus
  rev: 7f84f4286a52ec63260adeb8398ca7814ae19422
- path: golang.org/x/tools
  rev: 1330b289ad6d59313d86910fa35d5022cd871e7f
- path: github.com/spf13/cobra
  rev: 4b86c66ef25470e678a4d6a372711d7050344ccc
- path: github.com/kr/fs
  rev: 2788f0dbd16903de03cb8186e5c7d97b69ad387b
- path: github.com/inconshreveable/mousetrap
  rev: 76626ae9c91c4f2a10f34cad8ce83ea42c93bb75
- path: golang.org/x/net
  rev: db8e4de5b2d6653f66aea53094624468caad15d2
- path: github.com/spf13/pflag
  rev: f735fdff4ffb34299727eb2e3c9abab588742d41
- path: github.com/cpuguy83/go-md2man
  rev: 71acacd42f85e5e82f70a55327789582a5200a90
- path: gopkg.in/check.v1
  rev: 11d3bc7aa68e238947792f30573146a3231fc0f1
- path: golang.org/x/text
  rev: 6368131e2e9977b23aa20c631034cb98d65461a7
- path: github.com/russross/blackfriday
  rev: 8cec3a854e68dba10faabbe31c089abf4a3e57a6
- path: github.com/shurcooL/sanitized_anchor_name
  rev: 244f5ac324cb97e1987ef901a0081a77bfd8e845
```

### List Dependencies
If you want to scan your code to find out how many third party dependencies are
in your project run `govend list`. You can specify a path and output formats.

Here is an example of a path: `govend list some/project/path`
```bash
github.com/spf13/cobra
github.com/kr/fs
gopkg.in/yaml.v2
github.com/jackspirou/importsplus
golang.org/x/tools/go/vcs
```

**JSON**
`govend list -f json`
```bash
[
  "github.com/spf13/cobra",
  "github.com/kr/fs",
  "gopkg.in/yaml.v2",
  "github.com/jackspirou/importsplus",
  "golang.org/x/tools/go/vcs"
]%  
```

**YAML**
`govend list -f yml`
```bash
- github.com/spf13/cobra
- github.com/kr/fs
- gopkg.in/yaml.v2
- github.com/jackspirou/importsplus
- golang.org/x/tools/go/vcs
```
**XML**
`govend list -f xml`
```bash
<string>github.com/spf13/cobra</string>
<string>github.com/kr/fs</string>
<string>gopkg.in/yaml.v2</string>
<string>github.com/jackspirou/importsplus</string>
<string>golang.org/x/tools/go/vcs</string>%
```

Known Issues
============

Does `govend` work on Windows platforms?

> It does, we have tested it, but some bugs may exist.

Contributing
============

### Can I Contribute?

> Please do! Simply fork the code and send a pull request.
