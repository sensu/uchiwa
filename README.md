# Uchiwa

*Uchiwa* is simple dashboard for the Sensu monitoring framework, built with node.js.

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

## Debugging
You may start the dashboard with the following command in order to enable verbose mode: `NODE_ENV="development" node app.js`
