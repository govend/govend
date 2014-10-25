![Golang Gopher](./images/small-gopher.png) ![GoVend](./images/govend.png)

What Is GoVend
==============

The command `govend` takes yet another stab at easily integrating dependency management with golang projects. While many different solutions already exist to manage third party golang packages, `govend` tries a more simple approach.

How It Works
============

In my (limited) experience, `go get` has been very effective for downloading and adding golang packages into a local development `$GOPATH`. Yet, when using `go get` as a step in production deployments it has not been as nearly effective. In my (again limited) experience, depending on the OS build, network environment, and hosting provider `go get` might fail.

`govend` solves this problem by pulling golang dependencies into your project repo. By creating a `deps.json` file that lists your `go get` dependencies, running `govend` will copy those packages into your desired repository directory.

This is achieved by wrapping `go get` and only managing the movement of dependency files and packages. It is true that pulling files and packages from the `$GOPATH` into repos is not a difficult task, but if used correctly this simple feature can **go** a long way... :)

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

Using a different directory than `_vendor`.
===========================================

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

### Why will some packages like `github.com/someone/someproject` not get pushed up when I commit?

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
