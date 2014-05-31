
/**
 * Module dependencies.
 */
var express = require('express'),
  http = require('http'),
  path = require('path'),
  async = require('async'),
  moment = require('moment'),
  app = express();

server = http.createServer(app);
io = require('socket.io').listen(server);
io.set('log level', 1);

var config = require('./config.js')
var Sensu = require('./lib/sensu.js').Sensu;
var Stats  = require('./lib/stats.js').Stats;
var clients = {};
var stats = {};

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
process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0"

/**
 * Authentification
 */
if (config.uchiwa.user && config.uchiwa.pass){
  var basicAuth = express.basicAuth(function(username, password) {
    return (username == config.uchiwa.user && password == config.uchiwa.pass);
  }, 'Restrict area, please identify');
  app.all('*', basicAuth);
}

/**
 * Error handling
 * NODE_ENV=development node app.js
 */
if ('development' == process.env.NODE_ENV) {
  console.log('Debugging enabled.')
  app.use(express.errorHandler({showStack: true, dumpExceptions: true}));
  io.set('log level', 3);
}
app.use(function(err, req, res, next) {
  console.log(err);
  res.send(500);
})

var sensu = new Sensu(config.sensu);
var stats = new Stats(config.uchiwa, sensu);

var getStashes = function(callback){
  sensu.getStashes(function(err, result){
    sensu.stashes = (err) ? {} : result;
    if (!err) sensu.getTimestamp(sensu.stashes, "timestamp", "last_check", function(){});
    if (!err) sensu.buildStashes(function(){});
    callback(err);
  });
};

var getClients = function(callback){
  sensu.getClients(function(err, result){
    sensu.clients = (err) ? {} : result;
    if (!err) sensu.getTimestamp(sensu.clients, "timestamp", "last_check", function(err){});
    callback(err);
  });
};

var getEvents = function(callback){
  sensu.getEvents(function(err, result){
    sensu.events = (err) ? {} : result; 
    if (!err) {
     sensu.buildEvents(callback);
    } 
    else {
      callback(err);
    }
  });
};

var getChecks = function(callback){
  sensu.getChecks(function(err, result){
    sensu.checks = (err) ? {} : result;
    if (!err) sensu.buildChecks(function(){});
    callback(err);
  });
};

var getClient = function(data, callback){
  sensu.getClient(data.name, function(err, result){
    var client = (err) ? {} : result;
    if (!err) sensu.sortHistory(client.history, "check", "last_status", function(err){});
    if (!err) sensu.getTimestamp(client.history, "last_execution", "last_check", function(err){});
    callback(err, client);
  });
};

var pull = function(){
  async.waterfall([
    getStashes,
    getChecks,
    getClients,
    getEvents,
    function(callback){
      sensu.sortEvents(sensu.events, "name", "status", callback);
    },
    function(callback){
      sensu.sortClients(sensu.clients, sensu.events, callback);
    },
    function(callback){
      sensu.sortByKey(sensu.checks, "name", callback);
    },
    function(callback){
      sensu.buildClients(callback);
    },
    function(callback){
      stats.getDashboard(callback);
    }
  ], function(err){
    if (err){
      console.log(err);
      io.sockets.emit('messenger', {content: JSON.stringify({"type": "error", "content": err})});
    }
    else {
      io.sockets.emit('sensu', {content: JSON.stringify(sensu)});
      io.sockets.emit('stats', {content: JSON.stringify(stats.dashboard)});
    }
  });
};
// Perform a pull on start and every config.uchiwa.refresh milliseconds
setInterval(pull, config.uchiwa.refresh);
pull();
/**
 * Listen for events
 */
io.sockets.on('connection', function (socket) {
  // Keep track of active clients
  clients[socket.id] = socket;

  // Remove client on disconnection
  socket.on('disconnect', function () {
    delete clients[socket.id];
  });

  socket.on('get_client', function (data){
    getClient(data, function(err, result){
      if (err){
        return console.error("Fatal error! " + err);
      } else {
        //console.log(result);
        clients[socket.id].emit('client', {content: JSON.stringify(result)});
      }
    });
  });
  socket.on('get_sensu', function (data){
    clients[socket.id].emit('sensu', {content: JSON.stringify(sensu)});
  });
  socket.on('get_stats', function (data){
    clients[socket.id].emit('stats', {content: JSON.stringify(stats.dashboard)});
  });
  socket.on('create_stash', function (data){
    sensu.postStash(data, function(err, result){
      if(err){
        clients[socket.id].emit('messenger', {content: JSON.stringify({"type": "error", "content": "<strong>Error!</strong> The stash was not created. Reason: " + err})});
      }
      else {
        clients[socket.id].emit('messenger', {content: JSON.stringify({"type": "success", "content": "<strong>Success!</strong> The stash has been created."})});
      }
    });
  });
  socket.on('delete_stash', function (data){
    sensu.deleteStash(data, function(err){
      if(err){
        clients[socket.id].emit('messenger', {content: JSON.stringify({"type": "error", "content": "<strong>Error!</strong> The stash was not deleted. Reason: " + err})});
      }
      else {
        clients[socket.id].emit('messenger', {content: JSON.stringify({"type": "success", "content": "<strong>Success!</strong> The stash has been deleted."})});
      }
    });
  });
  socket.on('resolve_event', function (data){
    sensu.resolveEvent(data, function(err){
      if(err){
        clients[socket.id].emit('messenger', {content: JSON.stringify({"type": "error", "content": "<strong>Error!</strong> The check was not resolved. Reason: " + err})});
      }
      else {
        clients[socket.id].emit('messenger', {content: JSON.stringify({"type": "success", "content": "<strong>Success!</strong> The check has been resolved."})});
      }
    });
  });
});

/**
 * Routing
 */
app.get('/', function(req,res) {
   res.render('index.html');
});
app.get('/index', function (req, res) {
    res.render('index.html');
});
app.get('/checks', function(req,res) {
  res.render('checks.html');
});
app.get('/clients', function(req,res) {
  res.render('clients.html');
});
app.get('/stashes', function(req,res) {
  res.render('stashes.html');
});

/**
 * Start server
 */
server.listen(app.get('port'), function () {
  console.log('Express server listening on port ' + app.get('port'));
});