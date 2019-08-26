# Ion Connect

CLI tool for interacting with the Ion Channel.

[![Build Status](https://travis-ci.org/ion-channel/ion-connect.svg?branch=master)](https://travis-ci.org/ion-channel/ion-connect)

## Download and install

The latest build is available at [https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz](https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz). Download, open the folder for the operating system being used and put the file somewhere it can be run.

## License

[Apache Software License 2.0](LICENSE.txt)

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

## Building from source

This information is only applicable if you are making changes to `ion-connect`.

First thing you will need is the Golang environment setup. [The install process](https://golang.org/doc/install) pretty simple you can use [brew](http://brew.sh) to install go.  You will also need to set the $GOPATH environment variable and point it at your go workspace. [GNU make](https://www.gnu.org/software/make/) is also used for running commands.

`make install` will download the Go linter tool.

`make test` will run the tests.

`make build` will build `ion-connect`.