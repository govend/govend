govend [![Build Status](https://travis-ci.org/gophersaurus/govend.svg?branch=master)](https://travis-ci.org/gophersaurus/govend)
============================================================================================================================

`govend` leverages the `GO15VENDOREXPERIMENT` to vendor dependencies.

**it does:**
* vendor dependencies to the `nth` degree
* try to be flexible for compatibility with any project
* fill in vendored import paths for you with `govend imports`
* use the `vendor` directory as specified in golang version 1.5

**it does not:**
* wrap the `go` command
* create a new project for you
* change environment variables, including `$GOPATH`

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
- path: github.com/codegangsta/cli
  rev: 9b2bd2b3489748d4d0a204fa4eb2ee9e89e0ebc6
- path: github.com/jackspirou/importsplus
  rev: 7f84f4286a52ec63260adeb8398ca7814ae19422
- path: gopkg.in/yaml.v2
  rev: 49c95bdc21843256fb6c4e0d370a05f24a0bf213
- path: github.com/kr/fs
  rev: 2788f0dbd16903de03cb8186e5c7d97b69ad387b
- path: golang.org/x/tools
  rev: 0c09ff325ac41535a3d5fb6d539c32aca981bada
- path: golang.org/x/net
  rev: 84ba27dd5b2d8135e9da1395277f2c9333a2ffda
- path: gopkg.in/check.v1
  rev: 64131543e7896d5bcc6bd5a76287eb75ea96c673
- path: golang.org/x/text
  rev: d48eb58d199dc14dfaafefabf361feff840ee06c
```

### Scan Code
If you want to scan your code to find out how many third party dependencies are
present run `govend list`.  You can even specify a path and output formats.

Here is an example of a path: `govend scan some/path/dir`
```bash
github.com/codegangsta/cli
github.com/jackspirou/importsplus
gopkg.in/yaml.v2
github.com/kr/fs
golang.org/x/tools/go/vcs
```

**JSON**
`govend list -f json`
```bash
[
  "github.com/codegangsta/cli",
  "github.com/jackspirou/importsplus",
  "gopkg.in/yaml.v2",
  "github.com/kr/fs",
  "golang.org/x/tools/go/vcs"
]%
```

**YAML**
`govend list -f yml`
```bash
- github.com/codegangsta/cli
- github.com/jackspirou/importsplus
- gopkg.in/yaml.v2
- github.com/kr/fs
- golang.org/x/tools/go/vcs
```
**XML**
`govend list -f xml`
```bash
<string>github.com/codegangsta/cli</string>
<string>github.com/jackspirou/importsplus</string>
<string>gopkg.in/yaml.v2</string>
<string>github.com/kr/fs</string>
<string>golang.org/x/tools/go/vcs</string>%  
```

Known Issues
============

Does `govend` work on Windows platforms?

> I think so, but it should be tested.  Let me know what you find.

Contributing
============

### Can I Contribute?

> Please do! Simply fork the code and send a pull request.
