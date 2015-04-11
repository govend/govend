govend [![Build Status](https://travis-ci.org/gophersaurus/govend.svg?branch=master)](https://travis-ci.org/gophersaurus/govend)
============================================================================================================================

The CLI tool `govend` takes a stab at golang dependency management.  While many
projects already exist to manage external golang packages, `govend` follows a
philosophy of minimal interaction with the user, the project,
the host's environment configuration, and the `go` command.

Essentially this means `govend` tries to be good at one thing, vendoring dependencies.

**govend does not try to:**
* change any environment variables, including `$GOPATH` (let's not mess with your
  personal config)
* create a new golang project for you (govend should work with any project)
* wrap the `go` command (if `go` gets a major update, we shouldn't break you)
* dump the dependency manifest in the project root (we hide behind the curtain,
  not show off)
* make you maintain the dependency file `internal/_vendor/vendors.yml`
(lots of copy/paste)

**govend does try to:**
* be compatible with any project
* use the `internal` directory as specified in golang version 1.4
* rewrite all import paths
* rewrite golang import comments such as `// import "github.com/org/proj"`
* make the long import paths easy to maintain with `govend imports`


Demo
====

![govend demo](https://raw.githubusercontent.com/gophersaurus/govend/master/images/govend_demo.gif)

Install
=======

```bash
$ go get github.com/gophersaurus/govend

$ cd project/root/path

$ govend -verbose
```

Options
=======

```bash
→ govend -h

NAME:
   govend - A CLI tool for vendoring golang packages.

USAGE:
   govend [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   scan, s	Scans a go project for external package dependencies
   imports, i	Rewrites imports prioritizing the projects vendor directory
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose		print things as they happen
   --help, -h		show help
   --version, -v	print the version
```

### Vendor Code
To vendor code run `govend -verbose` or simply `govend` while in the root directory of your project.  This command will vendor all your dependecies and create a simple `internal/_vendor/vendors.yml` file to maintain versions.

`govend` can understand a wide range of senarios, so when in doubt run `govend`.

```bash
→ govend -verbose
identifying project paths...                    complete
scanning for external unvendored packages...    5 packages found
will generate manifest...                       internal/_vendor/vendors.yml
identifying repositories...                     complete
downloading packages...
 ↓ https://github.com/codegangsta/cli (latest)
 ↓ https://github.com/jackspirou/importsplus (latest)
 ↓ https://gopkg.in/yaml.v2 (latest)
 ↓ https://github.com/kr/fs (latest)
 ↓ https://go.googlesource.com/tools (latest)
vendoring packages...                           complete
writing vendors.yml manifest...                 complete
rewriting import paths...                       complete
```

### `internal/_vendor/vendors.yml`

The `vendors.yml` file simply contains an array of import paths and commit revisions.

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
- path: golang.org/x/tools/go/vcs
  rev: ""
```

### Rewrite Imports
`govend` will always rewrite imports, but `govend imports` comes in handy while writing code.  `govend imports` acts exactly like `goimports`, but will prioritize the `internal/_vendor` directory.  If you have ever vendored packages you have probably noticed that `goimports` first pulls from unvendored `$GOPATH` packages.  This can be annoying so give `govend imports -w` a shot.

### Scan Code
If you want to scan your code to find out how many third party dependencies are present run `govend scan`.  You can even specify a path and output formats.

Here is an example of a path: `govend scan some/path/dir`
```bash
github.com/codegangsta/cli
github.com/jackspirou/importsplus
gopkg.in/yaml.v2
github.com/kr/fs
golang.org/x/tools/go/vcs
```

**JSON**
`govend scan -fmt=json`
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
`govend scan -fmt=yml`
```bash
- github.com/codegangsta/cli
- github.com/jackspirou/importsplus
- gopkg.in/yaml.v2
- github.com/kr/fs
- golang.org/x/tools/go/vcs
```
**XML**
`govend scan -fmt=xml`
```bash
<string>github.com/codegangsta/cli</string>
<string>github.com/jackspirou/importsplus</string>
<string>gopkg.in/yaml.v2</string>
<string>github.com/kr/fs</string>
<string>golang.org/x/tools/go/vcs</string>%  
```

How It Works
============

`govend` works by running the following steps below:

 1. Identify all relative file paths necessary for the current project.
 2. Identify all types of packages currently present in the project.
 3. If the `vendors.yml` manifest file exists, load it in memory.
 4. Verify vendored packages and treat bad ones as unvendored packages.
 5. Identify package repositories and filter out repo subpackages.
 6. Download and vendor packages.
 7. Write the `vendors.yml` manifest file.
 8. Rewrite import paths.

A highlevel visualize is below:

![alt text](https://raw.githubusercontent.com/jackspirou/govend/ft-rewrite/images/govend_flow.png "govend flow")

Another dependency solution?
============================

If your looking for other dependency solutions, here is my list:

**Leading Projects**
* `go get`
* [Godeps](https://github.com/tools/godep)
* [nut](https://github.com/jingweno/nut)

**Many Others**
* [PackageManagementTools](https://github.com/golang/go/wiki/PackageManagementTools)

In my experience, `go get` has been very effective for downloading and adding golang packages into a local development `$GOPATH`. Yet, you should not be using `go get` as a step for production deployments.  If you are doing that, please stop and use `govend` or `godeps` or `nut` or any third party dependency manager.  I had experiences that made me fear when  depending on the OS build, network environment, and hosting provider to ensure `go get` would not fail.

> "go get is nice for playing around, but if your going to do something serious like deploy binaries to production, your deploy to production script shouldn't involve fetching some random dude's stuff on github. - Brad Fitzpatrick"
http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=36m40s

So we all agree that `go get` is a bad idea.  What about `godeps`?  `godeps` may be perfect for you.  Some really big projects use `godeps` and I admire the author of `godeps`, but it doesn't do quite what I want.  `godeps` edits your `$GOPATH` and also wraps the `go` command like so... `godeps go build` and I want to keep my tools seperate.  I don't want to rely on `godeps` not messing up `go`.  Just my opinion.

K - so how about `nut`?  `nut` is much closer to what I want.  `nut` felt the same way about changing the `$GOPATH` and wrapping the `go` command so the author avoided that.  Good job `nut`!  What I don't like, is that `nut` has options for creating a new golang project.  I think that is beyond the scope for what a dependency management tool should do for you.  Also I think the `Nut.toml` file is odd, but I'm sure people think my choice of a `yaml` file is odd.  Finally the `Nut.toml` has options for keep track of your project name, version, authors, and email addresses.  Im not saying those are not nice features, I just think they should be some other tools problem.

Fine then, what about project `X`?  K - I have officially exhausted all of my knowledge of different golang dependency management tools.  I did this to create what I wanted - but if there is better tool out there let me know!

Why Try To Solve This Problem?
=============================

I like to think that this project and others were inspired by talks at **GopherCon14** (I am too poor to attend) and the **GoTeam Google I/O Golang Fireside Chat 2013** (still to poor to attend).

You can watch them online just like I do at these links below! (yay internets)

-	Fireside Chat (part 1)
	http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=4m30s

-	Fireside Chat (part 2)
	http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=36m40s

-	GopherCon14 SoundCloud Best Practices for Production Environments
	http://www.youtube.com/watch?v=Y1-RLAl7iOI&feature=youtu.be&t=20m5s

Known Issues
============

Does `govend` work on Windows platforms?

> I have no idea.  I think so, but it should be tested.  Let me know what you find.

Why will some packages in `vendor` not get pushed up when I commit?

> Take a look at your `.gitignore` and `.gitignore_global` files. I had an issue where one of these files ignored `*.com` which would include most third party golang packages.
>
> For your `.gitignore_global` I recommend the file below, but maybe changes to this will be needed as well:

[.gitignore_global](https://gist.github.com/jackspirou/eb8bcf296136056fa88d)
```yaml
# Compiled source #
###################
*.class
*.dll
*.exe
*.o
*.so

# Packages #
############
# it's better to unpack these files and commit the raw source
# git has its own built in compression methods
*.7z
*.dmg
*.gz
*.iso
*.jar
*.rar
*.tar
*.zip

# Logs and databases #
######################
*.log
*.sql
*.sqlite

# OS generated files #
######################
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
```


Contributing
============

### Can I Contribute?

> Please do! I need all the help I can get :)
