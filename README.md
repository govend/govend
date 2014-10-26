![Golang Gopher](./images/small-gopher.png) ![GoVend](./images/govend.png)

GoVend
======

The command `govend` takes yet another stab at easily integrating dependency management with golang projects. While many different solutions already exist to manage third party golang packages, `govend` tries a more simple approach.

Install
=======

```bash
$ go get github.com/jackspirou/govend
```

How It Works
============

> The very short answer:

Create a `deps.json` file that lists your `go get` dependencies and `govend` will copy those packages into your project repository.

> The very long answer:

In my (limited) experience, `go get` has been very effective for downloading and adding golang packages into a local development `$GOPATH`. Yet, when using `go get` as a step in a script for production deployments it has not been as nearly effective. In my (again limited) experience, depending on the OS build, network environment, and hosting provider `go get` might fail.

`govend` solves this problem by pulling golang dependencies into your project repo. By creating a `deps.json` file that lists your `go get` dependencies, running `govend` will copy those packages into your desired repository directory.

This is achieved by wrapping `go get` and only managing the movement of dependency files and packages. It is true that pulling files and packages from the `$GOPATH` into repos is not a difficult task, but if used correctly this simple feature can **go** a long way... :)

As [Brad Fitzpatrick](https://github.com/bradfitz) (a member of the golang team) said...

> "go get is nice for playing around, but if you your going to do something serious like deploy binarys to production, your deploy to production script should shouldn't involve fetching some random dude's stuff on github."

That inspirational quote can be heard here:

[![IMAGE ALT TEXT HERE](http://img.youtube.com/vi/p9VUCp98ay4/0.jpg)](http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=36m40s) http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=36m40s

Example
=======

> First create a `deps.json` file at the root of your project/repo, much like a README.md.

```javascript
{
    "deps": [
        "github.com/gorilla/mux",
        "gopkg.in/mgo.v2"
    ]
}
```

> Next run the command `govend` from your project/repo root.

```bash
$ govend

 ↓ github.com/gorilla/mux
 ↓ gopkg.in/mgo.v2

Vending complete
```

> Lastly change your imports to use the `_vendor` directory.

```go
package example

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"_vendor/github.com/gorilla/mux"
	"_vendor/gopkg.in/mgo.v2"
)
```

> Now `go build`!

Using a different directory than `_vendor`
==========================================

So you might have your own golang file structure you really like... and maybe that means `_vendor` being called "_vendor" or being located at the root of your project is a no go. Well just add an extra line to your `deps.json` file. An example is below:

> A `deps.json` file specifying the vendor directory.

```javascript
{
    "dir": "./app/vendor",
    "deps": [
        "github.com/gorilla/mux",
        "gopkg.in/mgo.v2"
    ]
}
```

Why This Project
================

I like to think that project was directly inspired by some talks at **GopherCon14** (I was too poor to attend) and the **GoTeam Google I/O Golang Fireside Chat 2013** (also to poor to attend).

You can watch them online just like I did at these links below! (yay internets)

-	Fireside Chat (part 1) http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=4m30s

-	Fireside Chat (part 2) http://www.youtube.com/watch?v=p9VUCp98ay4&feature=youtu.be&t=36m40s

-	GopherCon14 SoundCloud Best Practices for Production Environments http://www.youtube.com/watch?v=Y1-RLAl7iOI&feature=youtu.be&t=20m5s

Known Issues
============

### Does `govend` work on all OS platforms?

> Honestly `govend` is very new and it was just created to work on Macs right now. I haven't tested it on windows or linux machines, but it should work in theory.

### Does `govend` work for all version control platforms?

> Honestly `govend` is focused on working with GIT because I use GIT.
>
> Right now `govend` actually prunes out `.git` and `.gitignore` files from third party packages so that your project repo doesn't behave unexpectedly.
>
> Because `govend` wraps around `go get` it should work with all `go get` supported versioning control software... in theory.

### Why will some packages in `_vendor` not get pushed up when I commit?

> Take a look at your `.gitignore` and `.gitignore_global` files. I had an issue where one of these files ignored `*.com` which would include most third party golang packages.
>
> For your `.gitignore_global` I recommend the file below, but maybe changes should be made to that as well:

```yml
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
