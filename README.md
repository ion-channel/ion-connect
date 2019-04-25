# Ion Connect

CLI tool for interacting with the Ion Channel Service/API

Master Branch
[![Build Status](https://travis-ci.org/ion-channel/ion-connect.svg?branch=master)](https://travis-ci.org/ion-channel/ion-connect)

## Latest Build

[https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz](https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz)

### Install on RHEL OS types with yum/dnf
Add the following to /etc/yum.repos.d/ion-channel.repo

```sh
[ion-channel]
name=Ion Channel Repo
baseurl=https://yum.fury.io/ionchannel/
enabled=1
gpgcheck=0
```

Then run:

```sh
sudo dnf clean all; sudo dnf install ion-connect
```

## License

[Apache Software License 2.0](LICENSE.txt)

## Let's build it!

First thing you will need is the Golang environment setup. [The install process](https://golang.org/doc/install) pretty simple you can use [brew](http://brew.sh) to install go.  You will also need to set the $GOPATH environment variable and point it at your go workspace.

```sh
$ brew install go
$ brew install glide
$ export GOPATH=$HOME/go
$ export GOBIN=$GOPATH/bin
```

Go provides a command for retrieving dependencies called `go get`.  Since ion-connect is currently private it may be helpful (and more generally helpful) to force git to use ssh by default with  `git config --global url."git@github.com:".insteadOf "https://github.com/"`

```sh
$ go get github.com/ion-channel/ion-connect
```

then grab our Go dependencies:

```sh
$ glide install
```

finally you can build ion-connect with the following:

```sh
$ go build github.com/ion-channel/ion-connect
```

## Don't forget the tests

NOTE: test dependencies are also managed with `glide install` as above.

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

### To get stuff 'cross-compiled':

```sh
$ GOOS=windows go build -o ion-connect/windows/bin/ion-connect.exe ./

$ GOOS=linux go build -o ion-connect/linux/bin/ion-connect ./

$ GOOS=darwin go build -o ion-connect/darwin/bin/ion-connect ./
```

## Your wish

Ion Connect provides a setup command called *configure*.  This should probably be the first command you run.  You will be prompted for your Ion Channel Secret Key which will be provided by an Ion Channel staff member. Contact us at <ion-connect@ionchannel.io>

```sh
$ ion-connect configure
Ion Channel Api Key []:
```

NOTE:  Ion Connect can also do configuration through environment variables.  The following variables are supported:

- `IONCHANNEL_SECRET_KEY` - allows the user to set the secret key used for authentication with Ion Channel
- `IONCHANNEL_ENDPOINT_URL` - allows the user to set the location of the Ion Channel api

The default endpoint url is the Ion Channel public/production environment.

You can then run various commands to query the Ion Channel system.  The best next step is the the help command.  From there you can see a list of top level commands and global options.

```sh
NAME:
   ion-connect - Interact with Ion Channel

USAGE:
   ion-connect [global flags] command [command flags] [arguments...]

VERSION:
   0.10.2

COMMANDS:
     scanner        set of commands for effecting artifacts or source code
     metadata       set of commands for parsing metadata from text
     ruleset        set of commands for managing rulesets
     analysis       set of commands for querying for projects analysis scan results
     project        set of commands for manipulating projects for your team
     dependency     set of commands for querying dependency data
     mail           set of commands for querying mail message data
     community      set of commands for querying a project's community data
     raw            set of commands for querying raw analysis and scan data
     vulnerability  set of commands for querying vulnerabilities
     configure      setup the Ion Channel secret key for later use
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        display debug logging
   --insecure     allow for insecure https connections
   --help, -h     show help
   --version, -v  print the version
```

Commands that are only supported in test by supplying the Ion Channel API endpoint url using an environment variable, similar to the following:

```
$ IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect metadata get-sentiment 'I love governance'
{
  "score": 0.925,
  "sentiment": "positive"
}
```

That's it! You are well on your way to world domination.
