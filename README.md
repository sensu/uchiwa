# Uchiwa

*Uchiwa* is a simple dashboard for the Sensu monitoring framework, built with node.js.

The dashboard is under active development, and major changes are not uncommon.

## Features

* Real-time updates with Socket.IO
* Client and checks stashes management
* Simple client details view
* Easy installation

## Screenshots

![Dashboard](http://palourde.github.io/images/uchiwa-dashboard.png)

![Dashboard](http://palourde.github.io/images/uchiwa-client.png)

## Installation

* Checkout the source: `git clone git@github.com:palourde/uchiwa.git`
* Install the dependencies: `npm install`

## Getting Started

* Copy **config.js.example** to **config.js** - modify your Sensu API information
* Start the dashboard: `node app.js`
* Browse your browser: `http://localhost:3000/`

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

## Debugging
You may start the dashboard with the following command in order to enable verbose mode: `NODE_ENV="development" node app.js`
