Govend [![Build Status](https://travis-ci.org/jackspirou/govend.svg?branch=master)](https://travis-ci.org/jackspirou/govend)
============================================================================================================================

The command `govend` takes yet another stab at solving golang dependency management. While many different solutions already exist to manage third party golang packages, `govend` tries a slightly different approach.

`govend` tries to be good at one thing, vendoring dependecies.

**govend does not try to:**
* change any enviorment variables, including `$GOPATH` your enviorment
* create a new golang project for you
* wrap the `go` command
* make you maintain the dependecy file `internal/_vendor/vendors.yml`
* dump the dependecy file in the root of your project

**govend does try to:**
* be compatible with any project
* use the `internal` directory as specified in golang version 1.4
* rewrite all import paths
* rewrite golang import comments such as `// import "github.com/org/proj"`
* make the long import paths easy to maintain with `govend imports`

Install
=======

```bash
$ go get github.com/jackspirou/govend

$ cd github.com/org/project-name

$ govend -verbose
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

K - so how about `nut`?  `nut` is much closer to what I want.  `nut` felt the same way about changing the `$GOPATH` and wrapping the `go` command so the author avoided that.  Good job `nut`!  What I don't like, is that `nut` has options for creating a new golang project.  I think that is beyond the scope for what a dependency managment tool should do for you.  Also I think the `Nut.toml` file is odd, but im sure people think my choice of a `yaml` file is odd.  Finally the `Nut.toml` has options for keep track of your project name, version, authors, and email addresses.  Im not saying those are not nice features, I just think they should be some other tools problem. 

Fine then, what about project `X`?  K - I have officially exahusted all of my knowledge of different golang dependency managment tools.  I did this to create what I wanted - but if there is better tool out there let me know! 

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
