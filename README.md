# Uchiwa

*Uchiwa* is a simple dashboard for the Sensu monitoring framework, built with node.js.

The dashboard is under active development, and major changes are not uncommon.

[![Code Climate](https://codeclimate.com/github/palourde/uchiwa.png)](https://codeclimate.com/github/palourde/uchiwa)

## Features

* Support of multiple Sensu APIs
* Real-time updates with Socket.IO
* Client and checks stashes management
* Easily filter events, clients, stashes and events
* Simple client details view
* Easy installation

## Screenshots

![Dashboard](http://palourde.github.io/images/uchiwa-dashboard.png)

![Dashboard](http://palourde.github.io/images/uchiwa-client.png)

## Installation

* Checkout the source: `git clone git@github.com:palourde/uchiwa.git`
* Install bower on your system: `npm install -g bower` 
* Install the dependencies: `npm install`

## Getting Started

* Copy **config.js.example** to **config.js** - modify your Sensu API information. See configuration section below
* Start the dashboard: `node app.js`
* Browse your browser: `http://localhost:3000/`

### Migration from 0.0.x to 0.1.x

With the support of mutiple Sensu APIs, the configuration structure has been modified. To configure multiple APIs, simply refer yourself to the **config.js.example** file.

Also make sure to run `npm install` to install any missing dependencies.

### Use nginx as proxy

The first thing you need is Nginx **1.3.13** or higher, since previous versions do not support websocket connections.

Then, you simply need to open up the Nginx configuration file and add the following route to your virtual server:
```
location / {
  proxy_pass http://localhost:3000;
  proxy_http_version 1.1;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  proxy_set_header Host $host;
}
```

In case you want the dashboard to be accessible within a certain path on the proxy, let's say /uchiwa, simply use the following block instead:
```
 location ~ (/uchiwa/|/socket.io/) {
  proxy_pass http://localhost:3000;
  proxy_http_version 1.1;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  proxy_set_header Host $host;

  rewrite /uchiwa/(.*) /$1 break;
}
```

## Configuration
### sensu
- `host` - String: The address of the Sensu API.
- `ssl` - Boolean: Determines whether or not the API use a SSL certificate.
- `port` - Integer: The port of the Sensu API. The default value is *4567*.
- `user` - String: The username of the Sensu API. Leave empty for none.
- `pass` - String: The password of the Sensu API. Leave empty for none.
- `path` - String: The path of the Sensu API. Leave empty in case of doubt.
- `timeout` - Integer: Timeout for the Sensu API, in milliseconds. The default value is *5000*.

### uchiwa
- `user` - String: The username of the Uchiwa dashboard. Leave empty for none.
- `pass` - String: The password of the Uchiwa dashboard. Leave empty for none.
- `stats` - Integer: Determines the retention, in minutes, of graphics data. The default value is *10*.
- `refresh` - Integer: Determines the interval to pull the Sensu API, in milliseconds. The default value is *10000*.

## Docker

This application comes pre-packaged in a docker container for easy deployment.

There are two ways of running this container:

### Docker with a config file.

Make a config.js file for the application, and then launch the uchiwa container with the config mounted as a volume.

    mkdir config
    cp ~/uchiwa/config.js.example config/config.js
    # Edit config
    docker run -v config:/config bobtfish:uchiwa
    # Docker will EXPOSE port 3000 by default, check where this is mapped on the host and browse to it.

### Docker with environment variables

You can instead use environment variables to configure the application. Host is fixed to 0.0.0.0 and port to 3000,
but the other settings can be set:

- `UCHIWA_USER`
- `UCHIWA_PASS`
- `UCHIWA_STATS`
- `UCHIWA_REFRESH`

And configuring an API is done with other environment variables which are designed to fit into Docker's
container links (allowing you to point uchiwa at an API just be --linking it to that container)

You can link multiple APIs by providing multiple sets of environment variables with different prefixes.

These variables are mandatory.

- `API1_PORT_4567_TCP_PORT` - The port for the API, usually 4567
- `API1_PORT_4567_TCP_ADDR` - The hostname or IP for the API

These variables are optional

- `API1_UCHIWA_NAME`
- `API1_UCHIWA_SSL`
- `API1_UCHIWA_USER`
- `API1_UCHIWA_PASS`
- `API1_UCHIWA_PATH`
- `API1_UCHIWA_TIMEOUT`

An example of starting the container with the minimum set of environment needed would be:

  docker run -i -t -p 3000 -e API1_PORT_4567_TCP_PORT=3000 -e API1_PORT_4567_TCP_ADDR="1.1.1.1" bobtfish/uchiwa

## Debugging
You may start the dashboard with the following command in order to enable verbose mode: `NODE_ENV="development" node app.js`

## Contributing
Everyone is welcome to submit patches. Whether your pull request is a bug fix or introduces new classes or functions to the project, we kindly ask that you include tests for your changes. Even if it's just a small improvement, a test is necessary to ensure the bug is never re-introduced.

### Testing
You should always make sure to have all dependencies installed (`npm install`)

#### Unit testing
Simply run `npm test`

#### E2E testing
1. Clone (this)[https://github.com/palourde/uchiwa-sensu] cookbook (`git clone git@github.com:palourde/uchiwa-sensu.git`)
2. Boot the virtual machines (`vagrant up`)
3. Use the following configuration file (*config.js*):
```
module.exports = {
  sensu: [
    {
      name: "0.12.6",
      host: '10.20.30.40',
      ssl: false,
      port: 4567,
      user: '',
      pass: '',
      path: '',
      timeout: 5000
    },
    {
      name: "0.13.0",
      host: '10.20.30.41',
      ssl: false,
      port: 4567,
      user: '',
      pass: '',
      path: '',
      timeout: 5000
    }
  ],
  uchiwa: {
    user: '',
    pass: '',
    retention: 10,
    refresh: 10000
  }
}
```
4. Run E2E tests (`npm run protractor`)

## Authors
* Author: [Simon Plourde][author] (<simon.plourde@gmail.com>)
* Contributor: Ethan Hann (<ethanhann@gmail.com>)

## License
Apache 2.0 (see [LICENSE][license])

[author]:                 https://github.com/palourde
[license]:                https://github.com/palourde/uchiwa/blob/master/LICENSE
