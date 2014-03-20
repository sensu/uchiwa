# uchiwa

A dashboard for Sensu.

## Install

* Checkout the source: `git clone git@github.com:palourde/uchiwa.git`
* Install the dependencies: `npm install`

## Getting Started

* Copy **config.js.example** to **config.js** - modify your Sensu API hostname
* Start the dashboard: `node app.js`
* Browse your browser: `http://localhost:3000/`

## Debugging
You may start the dashboard with the following command in order to enable verbose mode: `NODE_ENV="development" node app.js`
