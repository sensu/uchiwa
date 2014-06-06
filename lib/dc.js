var async = require('async');
var moment = require('moment');
var _ = require('underscore');
var Sensu = require('./sensu.js').Sensu;

function Dc(config) {
  this.name = config.name;
  this.sensu = new Sensu(config);
  this.style = "";
  this.criticals = 0;
  this.warnings = 0;
  this.events = 0;
  this.clients = 0;
  this.stashes = 0;
  this.checks = 0;
}

Dc.prototype.build = function(){
  this.criticals = (_.isEmpty(this.sensu.events)) ? 0 : this.sensu.events.filter(function (e){ return e.check.status === 2 }).length;
  this.warnings = (_.isEmpty(this.sensu.events)) ? 0 : this.sensu.events.filter(function (e){ return e.check.status === 1 }).length;
  this.clients = this.sensu.clients.length;
  this.events = this.sensu.events.length;
  this.stashes = this.sensu.stashes.length;
  this.checks = this.sensu.checks.length;
  this.style = (this.criticals > 0) ? "critical" : (this.warnings > 0) ? "warning" : "success";
};

Dc.prototype.getStashes = function(callback){
  var self = this;
  this.sensu.getStashes(function(err, result){
    self.sensu.stashes = (err) ? {} : result;
    if (!err) self.sensu.getTimestamp(self.sensu.stashes, "timestamp", "last_check", function(){});
    if (!err) self.sensu.buildStashes(function(){});
    callback(err);
  });
};

Dc.prototype.getClients = function(callback){
  var self = this;
  this.sensu.getClients(function(err, result){
    self.sensu.clients = (err) ? {} : result;
    if (!err) self.sensu.getTimestamp(self.sensu.clients, "timestamp", "last_check", function(err){});
    callback(err);
  });
};

Dc.prototype.getEvents = function(callback){
  var self = this;
  this.sensu.getEvents(function(err, result){
    self.sensu.events = (err) ? {} : result; 
    if (!err) {
     self.sensu.buildEvents(callback);
    } 
    else {
      callback(err);
    }
  });
};

Dc.prototype.getChecks = function(callback){
  var self = this;
  this.sensu.getChecks(function(err, result){
    self.sensu.checks = (err) ? {} : result;
    if (!err) self.sensu.buildChecks(function(){});
    callback(err);
  });
};

Dc.prototype.getClient = function(clientName, callback){
  var self = this;
  this.sensu.getClient(clientName, function(err, result){
    var client = (err) ? {} : result;
    if (!err) self.sensu.sortHistory(client.history, "check", "last_status", function(err){});
    if (!err) self.sensu.getTimestamp(client.history, "last_execution", "last_check", function(err){});
    callback(err, client);
  });
};

Dc.prototype.pull = function(next){
  var self = this;
  async.waterfall([
    self.getStashes.bind(this),
    self.getChecks.bind(this),
    self.getClients.bind(this),
    self.getEvents.bind(this),
    function(callback){
      self.sensu.sortEvents(self.sensu.events, "name", "status", callback);
    },
    function(callback){
      self.sensu.sortClients(self.sensu.clients, self.sensu.events, callback);
    },
    function(callback){
      self.sensu.sortByKey(self.sensu.checks, "name", callback);
    },
    function(callback){
      self.sensu.buildClients(callback);
    }
  ], function(err){
    if (err){
      console.log(moment().format('YYYY[-]MM[-]DD HH[:]mm[:]ss') + " [error] Sensu API " + self.name + " returned \"" + err + "\"");
      io.sockets.emit('messenger', {content: JSON.stringify({"type": "error", "content": "<strong>Error!</strong> Sensu API " + self.name + " returned \"" + err + "\""})});
      next();
    }
    else {
      next();
    }
  });
};

exports.Dc = Dc;