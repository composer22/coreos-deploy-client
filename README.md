# coreos-deploy-client
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/composer22/coreos-deploy-client.svg?branch=master)](http://travis-ci.org/composer22/coreos-deploy-client)
[![Current Release](https://img.shields.io/badge/release-v0.0.1-brightgreen.svg)](https://github.com/composer22/coreos-deploy-client/releases/tag/v0.0.1)
[![Coverage Status](https://coveralls.io/repos/composer22/coreos-deploy-client/badge.svg?branch=master)](https://coveralls.io/r/composer22/coreos-deploy-client?branch=master)

A CLI for interacting with coreos-deploy written in [Go.](http://golang.org)

## About

This client application can be used to submit deploy requests to coreos-deploy server
or check the status of a previous deploy request to the server.

## Requirements

You need to server source:

go get github.com/composer22/coreos-deploy

You will also need a valid API token setup there to access the server.

Etcd2 key values should be submitted as a text file with a space delimeter
between key and value:
```
/my/etcd2/key1 somevalue1
/my/etcd2/key2 somevalue2
/my/etcd2/key3 another value with embedded spaces.
```
## Usage

This command performs two functions:

* Submitting a deploy request to the coreos-deploy server in a CoreOS cluster.
* Retrieving the result of a previous deploy request.

When submitting a deploy request, a unique deploy ID is returned.  Use this in a
subsequent call to retrieve the result.

```
Description: coreos-deploy-client is a CLI for deploying services to the coreos-deploy API.

Usage: coreos-deploy-client [options...]

Server options:
    -n, --name NAME                  NAME of the service (mandatory).
    -r, --service_version VERSION    VERSION of the service (mandatory).
    -i, --instances INSTANCES        Number of INSTANCES to deploy. (default: 2).
    -t, --template_filepath TEMPLATE Path and filename to the unit .service TEMPLATE (mandatory).
    -e, --etcd2_filepath ETCD2FILE   Path and filename to the etcd2 key/value ETCD2FILE.
    -b, --bearer_token TOKEN         The API authorization TOKEN for the server.
    -u, --url URL                    URL of the coreos-deploy server.

    -p, --deploy_id ID               Lookup the status of a previous deployment.

    -d, --debug                      Enable debugging output (default: false)

Common options:
    -h, --help                       Show this message
    -V, --version                    Show version

Examples:

   # Deploy a service and return a deploy ID...
    coreos-deploy-client -n my-application -r 1.0.1 -i 2 \
	 -t /path/to/my-application@.service \
	 -e /path/to/my-application.etcd2 -b AP1T0K3N \
	 -u http://coreos-dev.example.com:80

	# Check the status of a recent deploy...
	coreos-deploy-client -b AP1T0K3N -u http://coreos-dev.example.com:80 \
	 -p DC8D9C2E-8161-4FC0-937F-4CA7037970D5

```
## CLI Bash Wrapper

A CLI bash script is provided for ease of use. The script prompts the user
for parameters and provides environment variables for the API token and
service URL.  To review help on the command, type one of the following:
```
./coreos-deploy-client.sh help
./coreos-deploy-client.sh help deploy
./coreos-deploy-client.sh help status

```
## Building

This code currently requires version 1.42 or higher of Go.

Information on Golang installation, including pre-built binaries, is available at
<http://golang.org/doc/install>.

Run `go version` to see the version of Go which you have installed.

Run `go build` inside the directory to build.

Run `go test ./...` to run the unit regression tests.

A successful build run produces no messages and creates an executable called `coreos-deploy-client` in this
directory.

Run `go help` for more guidance, and visit <http://golang.org/> for tutorials, presentations, references and more.

Run `./build.sh` for multiple platforms.

## License

(The MIT License)

Copyright (c) 2015 Pyxxel Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to
deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.
