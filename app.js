'use strict';

// Load Modules
var express = require('express');
var http = require('http');
var path = require('path');
var moment = require('moment');
var app = express();
var server = http.createServer(app);
var io = require('socket.io')(server);

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

// Express Configuration
app.set('config', config);
app.set('port', process.env.PORT || config.uchiwa.port);
app.set('host', process.env.HOST || config.uchiwa.host);
app.engine('.html', require('ejs').__express);
app.set('views', path.join(__dirname, 'public'));
app.set('view engine', 'html');
app.use(express.logger('dev'));
app.use(express.json());
app.use(express.urlencoded());
app.use(express.methodOverride());
app.use(app.router);
app.use(express.static(path.join(__dirname, 'public')));
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

/**
 * Error handling
 * DEBUG=* NODE_ENV=development node app.js
 */
if ('development' === process.env.NODE_ENV) {
  console.log('Debugging enabled.');
  app.use(express.errorHandler({showStack: true, dumpExceptions: true}));
}

/* jshint ignore:start */
app.use(function (err, req, res, next) {
  console.log(err);
  res.send(500);
});
/* jshint ignore:end */

// Authentification
if (config.uchiwa.user && config.uchiwa.pass) { app.all('*', authentication.basic); }

// Get Datacenters
config.sensu.forEach(function (configuration) {
  datacenters.push(new Dc(configuration));
});

// Pull & Push Sensu data
var refreshData = function () {
  pusher.pull(io, sensu, datacenters, function (result) {
    sensu = result;
    pusher.push(io, sensu, function () {});
  });
};
setInterval(refreshData, config.uchiwa.refresh);
refreshData();

// Listen for Socket.IO messages
io.on('connection', function (socket) { listeners.listen(socket, sensu, datacenters, publicConfig); });

// Status Page
app.get('/health/:component?', function(req, res){
  health.get(req, res, sensu, config);
});

// Start Server
server.listen(app.get('port'), app.get('host'), function () {
  console.log('Uchiwa is now listening on %s:%s', app.get('host'), app.get('port'));
});