'use strict';

/**
 * Module dependencies.
 */
var fs = require('fs');
var yargs = require('yargs')
    .describe('c', 'Load a config (file)')
    .alias('c', 'config')
    .alias('c', 'config-file')
    .alias('c', 'config_file')
    .default('c', './config.json');
var argv = yargs.argv;

if (!fs.existsSync(argv.c)) {
  yargs.showHelp();
  console.log('Config file must exist and be readable.');
  process.exit(1);
}

try {
  var config = require(argv.c);
} catch (e) {
  console.log('Syntax error with the config file ' + argv.c);
  process.exit(1);
}

var express = require('express'),
  http = require('http'),
  path = require('path'),
  async = require('async'),
  _ = require('underscore'),
  app = express(),
  server = http.createServer(app);

var io = require('socket.io').listen(server);
io.set('log level', 1);

var Dc = require('./lib/dc.js').Dc;
var Stats = require('./lib/stats.js').Stats;
var clients = {};
var stats = {};

/**
 * App configuration
 */
var port = config.uchiwa.port || 3000;
var host = config.uchiwa.host || '0.0.0.0';
app.set('port', process.env.PORT || port);
app.set('host', process.env.HOST || host);
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
 * Authentification
 */
if (config.uchiwa.user && config.uchiwa.pass) {
  var basicAuth = express.basicAuth(function (username, password) {
    return (username === config.uchiwa.user && password === config.uchiwa.pass);
  }, 'Restrict area, please identify');
  app.all('*', basicAuth);
}

/**
 * Error handling
 * NODE_ENV=development node app.js
 */
if ('development' === process.env.NODE_ENV) {
  console.log('Debugging enabled.');
  app.use(express.errorHandler({showStack: true, dumpExceptions: true}));
  io.set('log level', 3);
}
app.use(function (err, req, res, next) {
  console.log(err);
  res.send(500);
});

// Backward compatibility with uchiwa < 0.1.0
if (!_.isArray(config.sensu)) {
  config.sensu = [config.sensu];
  config.sensu[0].name = config.sensu[0].host;
}

var sensu = {};
var stats = new Stats(config.uchiwa);
var datacenters = [];
config.sensu.forEach(function (configuration) {
  datacenters.push(new Dc(configuration));
});

var pull = function () {
  var i = 0;
  sensu = {checks: [], clients: [], dc: [], events: [], stashes: []};
  async.eachSeries(datacenters, function (datacenter, nextDc) {
    datacenter.pull(function () {
        var aggregate = function (callback) {
          var attributes = ['checks', 'clients', 'events', 'stashes'];
          async.each(attributes, function (attribute, nextAttribute) {
            sensu[attribute][i] = [];
            async.each(datacenter.sensu[attribute], function (item, nextItem) {
              item.dc = datacenter.name;
              sensu[attribute][i].push(item);
              nextItem();
            }, function () {
              nextAttribute();
            });
          }, function () {
            callback();
          });
        };
        aggregate(function () {
          i++;
          datacenter.build();
          sensu.dc.push({
            name: datacenter.name,
            style: datacenter.style,
            clients: datacenter.clients,
            events: datacenter.events,
            stashes: datacenter.stashes,
            checks: datacenter.checks
          });
          nextDc();
        });
      },
      function (messageContent) {
        io.sockets.emit('messenger', {
          content: messageContent
        });
      }
    );
  }, function () {
    io.sockets.emit('sensu', {content: JSON.stringify(sensu)});

    // Update stats
    stats.getDashboard(sensu);
    io.sockets.emit('stats', {content: JSON.stringify(stats.dashboard)});
  });
};

// Perform a pull on start and every config.uchiwa.refresh milliseconds
setInterval(pull, config.uchiwa.refresh);
pull();

// Return DC object and check client if any specified
var getDc = function (data, callback) {
  if (datacenters.length === 0) {
    return callback('<strong>Error!</strong> No datacenters found.');
  }
  var dc = datacenters.filter(function (e) {
    return e.name === data.dc;
  });
  if (dc.length !== 1) {
    return callback('<strong>Error!</strong> The datacenter ' + data.dc + ' was not found.');
  }
  if (_.has(data, 'client')) {
    if (dc[0].sensu.clients.length === 0) {
      return callback('<strong>Error!</strong> No clients found.');
    }
    var client = dc[0].sensu.clients.filter(function (e) {
      return e.name === data.client;
    });
    if (client.length !== 1) {
      return callback('<strong>Error!</strong> The client ' + data.client + ' was not found.');
    }
  }
  callback(null, dc[0]);
};


/**
 * Listen for Socket.IO messages
 */
io.sockets.on('connection', function (socket) {
  // Keep track of active clients
  clients[socket.id] = socket;

  // Remove client on disconnection
  socket.on('disconnect', function () {
    delete clients[socket.id];
  });

  socket.on('get_sensu', function () {
    clients[socket.id].emit('sensu', {content: JSON.stringify(sensu)});
  });

  socket.on('get_stats', function () {
    clients[socket.id].emit('stats', {content: JSON.stringify(stats.dashboard)});
  });

  socket.on('get_client', function (data) {
    getDc(data, function (err, result) {
      if (err) {
        clients[socket.id].emit('messenger', {content: JSON.stringify({'type': 'error', 'content': err})});
      }
      else {
        result.getClient(data.client, function (err, result) {
          if (err) {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'error',
                'content': '<strong>Error!</strong> ' + err
              })
            });
          }
          else {
            clients[socket.id].emit('client', {content: JSON.stringify(result)});
          }
        });
      }
    });
  });

  socket.on('delete_client', function (data) {
    data = JSON.parse(data);
    getDc(data, function (err, result) {
      if (err) {
        clients[socket.id].emit('messenger', {content: JSON.stringify({'type': 'error', 'content': err})});
      }
      else {
        result.sensu.delete('clients', data.payload, function (err) {
          if (err) {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'error',
                'content': '<strong>Error!</strong> The client was not deleted. Reason: ' + err
              })
            });
          }
          else {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'success',
                'content': '<strong>Success!</strong> The client has been deleted.'
              })
            });
          }
        });
      }
    });
  });

  socket.on('create_stash', function (data) {
    data = JSON.parse(data);
    getDc(data, function (err, result) {
      if (err) {
        clients[socket.id].emit('messenger', {content: JSON.stringify({'type': 'error', 'content': err})});
      }
      else {
        result.sensu.post('stashes', JSON.stringify(data.payload), function (err) {
          if (err) {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'error',
                'content': '<strong>Error!</strong> The stash was not created. Reason: ' + err
              })
            });
          }
          else {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'success',
                'content': '<strong>Success!</strong> The stash has been created.'
              })
            });
          }
        });
      }
    });
  });

  socket.on('delete_stash', function (data) {
    data = JSON.parse(data);
    getDc(data, function (err, result) {
      if (err) {
        clients[socket.id].emit('messenger', {content: JSON.stringify({'type': 'error', 'content': err})});
      }
      else {
        result.sensu.delete('stashes', data.payload, function (err) {
          if (err) {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'error',
                'content': '<strong>Error!</strong> The stash was not deleted. Reason: ' + err
              })
            });
          }
          else {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'success',
                'content': '<strong>Success!</strong> The stash has been deleted.'
              })
            });
          }
        });
      }
    });
  });

  socket.on('resolve_event', function (data) {
    data = JSON.parse(data);
    getDc(data, function (err, result) {
      if (err) {
        clients[socket.id].emit('messenger', {content: JSON.stringify({'type': 'error', 'content': err})});
      }
      else {
        result.sensu.post('resolve', JSON.stringify(data.payload), function (err) {
          if (err) {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'error',
                'content': '<strong>Error!</strong> The check was not resolved. Reason: ' + err
              })
            });
          }
          else {
            clients[socket.id].emit('messenger', {
              content: JSON.stringify({
                'type': 'success',
                'content': '<strong>Success!</strong> The check has been resolved.'
              })
            });
          }
        });
      }
    });
  });

});


/**
 * Start server
 */
server.listen(app.get('port'), app.get('host'), function () {
  console.log('Uchiwa is now listening on %s:%s', app.get('host'), app.get('port'));
});
