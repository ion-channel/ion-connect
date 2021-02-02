# Ion Connect
[![Build Status](https://travis-ci.com/ion-channel/ion-connect.svg?branch=master)](https://travis-ci.com/ion-channel/ion-connect)
[![Go Report Card](https://goreportcard.com/badge/github.com/ion-channel/ion-connect)](https://goreportcard.com/report/github.com/ion-channel/ion-connect)
[![GoDoc](https://godoc.org/github.com/ion-channel/ion-connect?status.svg)](https://godoc.org/github.com/ion-channel/ion-connect)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/ion-channel/ion-connect/blob/master/LICENSE.md)
[![Release](https://img.shields.io/github/release/ion-channel/ion-connect.svg)](https://github.com/ion-channel/ion-connect/releases/latest)

CLI tool for interacting with the Ion Channel.

# Requirements
Golang Version 1.13 or higher

**and**

Ion API key

# Installation
Ion Connect can be installed from source, a binary download, or indirectly within a Docker
container.

## Go
```
go get github.com/ion-channel/ion-connect
cd $GOPATH/src/github.com/ion-channel/ion-connect
go build
go install
```

## Binary
[http://github.com/ion-channel/ion-connect/releases/latest](http://github.com/ion-channel/ion-connect/releases/latest)

## Docker
```
docker pull ionchannel/ion-connect
```

# Usage

## Configure

Ion Connect can be configured either by config file or through environment variables.

Config File:

```
ion-connect configure
```

Environment Variables:

```
export IONCHANNEL_SECRET_KEY=<your_api_key>
```

## Other Environments
The default environment for Ion Connect is the Ion Channel public/production environment. If for some reason you need to point Ion Connect to a different environment, you can use environment variables to change the API endpoint used.

```
export IONCHANNEL_ENDPOINT_URL=https://otherionurl
```

## Available Comands
Utilize the built in help to see what commands are available to Ion Connect.

```
ion-connect help
```

# Versioning
The project will be versioned in accordance with [Semver 2.0.0](http://semver.org).  See the [releases](https://github.com/ion-channel/ionic/releases) section for the latest version.  Until version 1.0.0 the project is considered to be unstable.

# License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).  See [LICENSE.md](./LICENSE.md) for more information.

## Updating dependencies with go modules
To update ionic:
`go get -u github.com/ion-channel/ionic@master && go mod vendor`

To update the linter requirements (to fix inconsistent vendoring in build):
`go get -u golang.org/x/lint && go get -u golang.org/x/tools && go mod vendor`
