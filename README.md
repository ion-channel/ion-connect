# Ion Connect

CLI tool for interacting with the Ion Channel.

[![Build Status](https://travis-ci.org/ion-channel/ion-connect.svg?branch=master)](https://travis-ci.org/ion-channel/ion-connect)

## Download and install

The latest build is available at [https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz](https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz). Download, open the folder for the operating system being used and put the file somewhere it can be run.

## License

[Apache Software License 2.0](LICENSE.txt)

## Your wish

Ion Connect provides a setup command called *configure*.  This should probably be the first command you run.  You will be prompted for your Ion Channel Secret Key which will be provided by an Ion Channel staff member. Contact us at <info@ionchannel.io>

```sh
$ ion-connect configure
Ion Channel Api Key []:
```

NOTE:  Ion Connect can also do configuration through environment variables.  The following variables are supported:

- `IONCHANNEL_SECRET_KEY` - allows the user to set the secret key used for authentication with Ion Channel
- `IONCHANNEL_ENDPOINT_URL` - allows the user to set the location of the Ion Channel api

The default endpoint url is the Ion Channel public/production environment.

You can then run various commands to query the Ion Channel system.  The best next step is the the help command to see what commands are available:

```sh
$ ion-connect help
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