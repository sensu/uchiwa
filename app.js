
/**
 * Module dependencies.
 */
var express = require('express'),
  http = require('http'),
  path = require('path'),
  async = require('async'),
  app = module.exports = express();

server = http.createServer(app);
io = require('socket.io').listen(server);
io.set('log level', 1);

var config = require('./config.js')
var Sensu = require('./lib/sensu.js').Sensu;

/**
 * App configuration
 */
app.set('port', process.env.PORT || 3000);
app.engine('.html', require('ejs').__express);
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'html');
app.use(express.favicon());
app.use(express.logger('dev'));
app.use(express.json());
app.use(express.urlencoded());
app.use(express.methodOverride());
app.use(app.router);
app.use(express.static(path.join(__dirname, 'public')));

/**
 * Error handling
 * NODE_ENV=development node app.js
 */
if ('development' == app.get('env')) {
  console.log('Debugging enabled.')
  app.use(express.errorHandler({showStack: true, dumpExceptions: true}));
  io.set('log level', 3);
}
app.use(function(err, req, res, next) {
  console.log(err);
  res.send(500);
})

var sensu = new Sensu(config.sensu);

var pull = function(){
  async.series([
    function(callback){
      sensu.getClients(function(err, result){
        sensu.clients = result;
        callback(err);
      });
    },
    function(callback){
      sensu.getTimestamp(sensu.clients, "timestamp", function(err){
        callback();
      });
    },
    function(callback){
      sensu.getEvents(function(err, result){
        sensu.events = result;
        callback(err);
      });
    },
    function(callback){
      sensu.getTimestamp(sensu.events, "issued", function(err){
        callback();
      });
    },
    function(callback){
      sensu.sortEvents(sensu.events, 'status', function(err){
        callback();
      });
    },
    function(callback){
      sensu.sortClients(sensu.clients, sensu.events, function(err){
        callback();
      });
    },
    function(callback){
      sensu.getChecks(function(err, result){
        sensu.checks = result;
        callback(err);
      });
    },
    function(callback){
      sensu.getStashes(function(err, result){
        sensu.stashes = result;
        callback(err);
      });
    },
    function(callback){
      sensu.getTimestamp(sensu.stashes, "content.timestamp", function(err){
        callback();
      });
    }
  ], function(err){
    if (!err){
      io.sockets.emit('events', {content: JSON.stringify(sensu.events)});
      io.sockets.emit('clients', {content: JSON.stringify(sensu.clients)});
      io.sockets.emit('checks', {content: JSON.stringify(sensu.checks)});
      io.sockets.emit('stashes', {content: JSON.stringify(sensu.stashes)});
    }
  });
};

// Perform a pull on start and every config.uchiwa.refresh milliseconds
pull();
setInterval(pull, config.uchiwa.refresh);

// Listen for events
io.sockets.on('connection', function (socket) {
  socket.on('get_client', function (data){
    async.series([
      function(callback){
        sensu.getClient(data.name, function(err, result){
          sensu.client = result;
          callback(err);
        });
      },
      function(callback){
        sensu.sortEvents(sensu.client, "last_status", function(err){
          callback();
        });
      },
      function(callback){
        sensu.getTimestamp(sensu.client, "last_execution", function(err){
          callback();
        });
      }
    ], function(err){
      if (err) return console.error("Fatal error! " + err);
      io.sockets.emit('client', {content: JSON.stringify(sensu.client)});
    });
  });
  socket.on('create_stash', function (data){
    sensu.postStash(data, function(err, result){
      if(err){
        io.sockets.emit('messenger', {content: JSON.stringify({"type": "error", "page": "client-details", "content": "<strong>Error!</strong> The stash was not created. Reason: " + err})});
      }
      else {
        console.log("success");
        io.sockets.emit('messenger', {content: JSON.stringify({"type": "success", "page": "client-details", "content": "<strong>Success!</strong> The stash has been created."})});
      }
    });
  });
  socket.on('delete_stash', function (data){
    sensu.deleteStash(data, function(err){
      if(err){
        io.sockets.emit('messenger', {content: JSON.stringify({"type": "error", "page": "client-details", "content": "<strong>Error!</strong> The stash was not deleted. Reason: " + err})});
      }
      else {
        io.sockets.emit('messenger', {content: JSON.stringify({"type": "success", "page": "client-details", "content": "<strong>Success!</strong> The stash has been deleted."})});
      }
    });
  });
});

/**
 * Routing
 */
app.get('/', function(req,res) {
   res.render('index.html', {title: 'Index'});
});
app.get('/clients', function(req,res) {
  res.render('clients.html', {title: 'Clients'});
  io.sockets.on('connection', function (socket) {
    io.sockets.emit('events', {content: JSON.stringify(sensu.events)});
    io.sockets.emit('clients', {content: JSON.stringify(sensu.clients)});
    io.sockets.emit('checks', {content: JSON.stringify(sensu.checks)});
    io.sockets.emit('stashes', {content: JSON.stringify(sensu.stashes)});
  });
});
app.get('/events',function(req,res) {
  res.render('events.html', {title: 'Events'});
  io.sockets.on('connection', function (socket) {
    io.sockets.emit('events', {content: JSON.stringify(sensu.events)});
  });
});

/**
 * Start server
 */
server.listen(app.get('port'), function () {
  console.log('Express server listening on port ' + app.get('port'));
});