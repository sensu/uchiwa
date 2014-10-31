# Uchiwa

*Uchiwa* is a simple dashboard for the Sensu monitoring framework, built with Go and AngularJS.

The dashboard is under active development, and major changes are not uncommon.

[![Build Status](https://travis-ci.org/sensu/uchiwa.svg?branch=master)](https://travis-ci.org/sensu/uchiwa)
[![Code   Climate](https://codeclimate.com/github/sensu/uchiwa/badges/gpa.svg)](https://codeclimate.com/github/sensu/uchiwa)

## Features

* Support of multiple Sensu APIs
* Real-time updates
* Client and checks stashes management
* Easily filter events, clients, stashes and events
* Simple client details view
* Easy installation

## Screenshots

![Dashboard](http://palourde.github.io/images/uchiwa-dashboard.png)

![Dashboard](http://palourde.github.io/images/uchiwa-client.png)

## Getting Started

### With packages

##### Using Sensu repositories
See [Sensu documentation](http://sensuapp.org/docs/0.13/dashboards_uchiwa)

### From source

**Prerequisites**
* Recent version of **git**
* Recent version of **npm** ([guide](https://github.com/joyent/node/wiki/installing-node.js-via-package-manager))
* Recent version of **go** ([guide](https://golang.org/doc/install))

**Installation**
* Checkout the source: `go get github.com/sensu/uchiwa && cd $GOPATH/src/github.com/sensu/uchiwa`
* Install third-party libraries:
  * With standard user: `npm install --production`
  * With root user: `npm install --production --unsafe-perm`
* Copy **config.json.example** to **config.json** - modify your Sensu API information. See configuration section below
* Start the dashboard: `go run uchiwa.go`
* Open your browser: `http://localhost:3000/`


## Configuration
### sensu
- `name` - String: Name of the Sensu API (used as datacenter name). If empty, a random one will be generated.
- `host` - String: The address of the Sensu API. **Required**.
- `port` - Integer: The port of the Sensu API. The default value is *4567*. **Required**
- `ssl` - Boolean: Determines whether or not to use the *HTTPS* protocol. The default value is *false*.
- `path` - String: The path of the Sensu API. Leave empty in case of doubt
- `user` - String: The username of the Sensu API. Leave empty for none.
- `pass` - String: The password of the Sensu API. Leave empty for none.
- `timeout` - Integer: Timeout for the Sensu API, in seconds. The default value is 5.

### uchiwa
- `host` - String: The address on which Uchiwa will listen. The default value is *0.0.0.0*.
- `port` - Integer: The port on which Uchiwa will listen. The default value is *3000*.
- `user` - String: The username of the Uchiwa dashboard. Leave empty for none.
- `pass` - String: The password of the Uchiwa dashboard. Leave empty for none.
- `refresh` - Integer: Determines the interval to pull the Sensu APIs, in seconds. The default value is *5*.

## Docker

This application comes pre-packaged in a docker container for easy deployment.

Make a config.json file for the application, and then launch the uchiwa container with the config mounted as a volume.

    # Create a folder that will be mount as a volume to the Docker container
    mkdir ~/uchiwa-config
    # Copy your uchiwa config into this last folder
    cp ~/uchiwa/config.json ~/uchiwa-config/config.json
    # Start Docker container. It will listen on port 3000 by default
    docker run -d -p 3000:3000 -v ~/uchiwa-config:/config uchiwa/uchiwa

## Health
You may easily monitor Uchiwa and the Sensu API endpoints with the **/health** page.

### /health
Returns Uchiwa and Sensu API status.
* success: 200
  * content: `{"uchiwa":"ok","sensu":{"0.12.6":{"output":"ok"},"0.13.0":{"output":"ok"}}}`
* error: 500
  * content: `{"uchiwa":"ok","sensu":{"0.12.6":{"output":"connect ECONNREFUSED"},"0.13.0":{"output":"ok"}}}`

### /health/uchiwa
Returns Uchiwa status.
* success: 200
  * content: `"ok"`
* error: 500
  * content: `"error"`

### /health/sensu
Returns Sensu API status.
* success: 200
  * content: `"{0.12.6":{"output":"ok"},"0.13.0":{"output":"ok"}}`
* error: 500
  * content: `{"0.12.6":{"output":"connect ECONNREFUSED"},"0.13.0":{"output":"ok"}}`

## Contributing
Everyone is welcome to submit patches. Whether your pull request is a bug fix or introduces new classes or functions to the project, we kindly ask that you include tests for your changes. Even if it's just a small improvement, a test is necessary to ensure the bug is never re-introduced.

### Testing
In order to install all the tools, please run `npm install`.

##### Backend (go)
The command `go test -v ./...` will execute the proper unit tests.

##### Frontend (AngularJS)
The command `grunt` will execute the proper linting and unit tests.

## Authors
* Author: [Simon Plourde][author] (<simon.plourde@gmail.com>)
* Contributor: [Justin Kolberg][amdprophet]
* Contributor: [ayan4m1][ayan4m1]
* Contributor: [Ethan Hann][ethanhann] (<ethanhann@gmail.com>)


## License
MIT (see [LICENSE][license])

[author]:                 https://github.com/palourde
[license]:                https://github.com/palourde/uchiwa/blob/master/LICENSE
[ethanhann]:              http://www.ethanhann.com/
[ayan4m1]:                https://github.com/ayan4m1
[amdprophet]:             https://github.com/amdprophet
