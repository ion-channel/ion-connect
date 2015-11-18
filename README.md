# Ion Connect

CLI tool for interacting with the Ion Channel API

[![Build Status](https://magnum.travis-ci.com/ion-channel/ion-connect.svg?token=AGRFpUr1LzvrKJ1SmsR3)](https://magnum.travis-ci.com/ion-channel/ion-connect)


## Let's build it!

First thing you will need is the Golang environment setup. [The install process](https://golang.org/doc/install) pretty simple you can use [brew](http://brew.sh) to install go.  You will also need to set the $GOPATH environment variable and point it at your go workspace.

```sh
$ brew install go
$ export $GOPATH=$HOME/go
```

Go provides a command for retrieving dependecies called `go get`.  Since ion-connect is currently private it may be helpful (and more generally helpful) to force git to use ssh by default with  `git config --global url."git@github.com:".insteadOf "https://github.com/"`

```sh
$ go get github.com/ion-channel/ion-connect
```

finally you can build ion-connect with the following:

```sh
$ go build github.com/ion-channel/ion-connect
```

## Don't forget the tests

Running the tests requires the install of a couple of test environment dependecies. You can use go to install these as well.

```sh
$ go get github.com/onsi/gomega
$ go get github.com/onsi/ginkgo
```

You can then run the tests with the following command:

```sh
$ go test -v github.com/ion-channel/ion-connect/...
```

## Install from source

Once you feel your tests are up to snuff you can use the go install command to install the ion-connect binary in the go bin directory.

```sh
$ go install github.com/ion-channel/ion-connect
```

If you've added the $GOHOME/bin to you path you should now be able to get ion-connect action going.

## Your wish

Ion Connect provides a setup command called *configure*.  This should probably be the first command you run.  You will be prompted for your Ion Channel Secret Key which will be provided by an Ion Channel staff member.

```sh
$ ion-connect configure
Ion Channel Api Key []:
```

You can then run various commands to query the Ion Channel system.  The best next step is the the help command.  From there you can see a list of top level commands and global options.

```sh
$ ion-connect help
NAME:
   ionconnect - Control AWS profiles

USAGE:
   ion-connect [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
   scanner	set of commands for effecting artifacts or source code
   configure	setup the Ion Channel secret key for later use
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug		display debug logging
   --help, -h		show help
   --version, -v	print the version
```

That's it! You are well on your way to world domination.
