'use strict';

// Express
var express = require('express.io');
var app = express();
app.http().io();

// Load Modules
var path = require('path');
var moment = require('moment');
var bunyan = require('bunyan');
var log = bunyan.createLogger({name: 'uchiwa', src: true});

// Uchiwa Librairies
var authentication = require('./lib/authentication.js');
var configuration = require('./lib/configuration.js');
var Dc = require('./lib/dc.js').Dc;
var listeners = require('./lib/listeners.js');
var pusher = require('./lib/pusher.js');
var health = require('./lib/health.js');

// Uchiwa Configuration
var sensu = {};
var datacenters = [];
var config = {};
configuration.get(function (result) { config = result; });
var publicConfig = configuration.public(config);
moment.defaultFormat = config.uchiwa.dateFormat;

// Authentification
app.set('config', config);

if (config.uchiwa.user && config.uchiwa.pass) { app.all('*', authentication.basic); }

// Express Configuration
app.set('port', process.env.PORT || config.uchiwa.port);
app.set('host', process.env.HOST || config.uchiwa.host);
app.engine('.html', require('ejs').__express);
app.set('views', path.join(__dirname, 'public'));
app.set('view engine', 'html');

app.use(express.static(path.join(__dirname, 'public')));

app.use(require('express-bunyan-logger')());
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

/**
 * Error handling
 * DEBUG=* NODE_ENV=development node app.js
 */
if ('development' === process.env.NODE_ENV) {
  log.info('Debugging enabled.');
  app.use(express.errorHandler({showStack: true, dumpExceptions: true}));
}

/* jshint ignore:start */
app.use(function (err, req, res, next) {
  log.error(err);
  res.send(500);
  next();
});
/* jshint ignore:end */

// Get Datacenters
config.sensu.forEach(function (configuration) {
  datacenters.push(new Dc(configuration));
});

// Pull & Push Sensu data
var refreshData = function () {
  pusher.pull(app, sensu, datacenters, function (result) {
    sensu = result;
    pusher.push(app, result, function () {});
  });
};
setInterval(refreshData, config.uchiwa.refresh);
refreshData();

// Listen for Socket.IO messages
listeners.listen(app, sensu, datacenters, publicConfig);

// Status Page
app.get('/health/:component?', function(req, res){
  health.get(req, res, sensu, config);
});

// Start Server
app.listen(app.get('port'), app.get('host'), function () {
  log.info('Uchiwa is now listening on %s:%s', app.get('host'), app.get('port'));
});