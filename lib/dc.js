'use strict';

var async = require('async');
var moment = require('moment');
var _ = require('underscore');
var Sensu = require('./sensu.js').Sensu;

function Dc(config) {
  this.name = config.name;
  this.sensu = new Sensu(config);
  this.style = '';
  this.criticals = 0;
  this.warnings = 0;
  this.events = 0;
  this.clients = 0;
  this.stashes = 0;
  this.checks = 0;
  this.info = {};
  this.health = 'ok';
}

Dc.prototype.build = function () {
  var self = this;
  var count = function (status) {
    if (_.isEmpty(self.sensu.events)) {
      return 0;
    }
    return self.sensu.events.filter(function (e) {
      return e.check.status === status;
    }).length;
  };
  this.criticals = count(2);
  this.warnings = count(1);
  this.unknown = count(3);
  this.clients = (_.isEmpty(this.sensu.clients)) ? 0 : this.sensu.clients.length;
  this.events = (_.isEmpty(this.sensu.events)) ? 0 : this.sensu.events.length;
  this.stashes = (_.isEmpty(this.sensu.stashes)) ? 0 : this.sensu.stashes.length;
  this.checks = (_.isEmpty(this.sensu.checks)) ? 0 : this.sensu.checks.length;
  this.style = (this.criticals > 0) ? 'critical' : (this.warnings > 0) ? 'warning' : (this.unknown > 0) ? 'unknown' : 'success';
};

Dc.prototype.getStashes = function (callback) {
  var self = this;
  this.sensu.get('stashes', function (err, result) {
    self.sensu.stashes = (err) ? {} : result;
    if (!err) {
      self.sensu.getTimestamp(self.sensu.stashes, 'timestamp', 'timestampParsed', function () {
      });
    }
    if (!err) {
      self.sensu.buildStashes(function () {
      });
    }
    callback(err);
  });
};

Dc.prototype.getClients = function (callback) {
  var self = this;
  this.sensu.get('clients', function (err, result) {
    self.sensu.clients = (err) ? {} : result;
    if (!err) {
      self.sensu.getTimestamp(self.sensu.clients, 'timestamp', 'lastCheck', function () {
      });
    }
    callback(err);
  });
};

Dc.prototype.getEvents = function (callback) {
  var self = this;
  this.sensu.get('events', function (err, result) {
    self.sensu.events = (err) ? {} : result;
    if (!err) {
      self.sensu.buildEvents(callback);
    }
    else {
      callback(err);
    }
  });
};

Dc.prototype.getChecks = function (callback) {
  var self = this;
  this.sensu.get('checks', function (err, result) {
    self.sensu.checks = (err) ? {} : result;
    if (!err) {
      self.sensu.buildChecks(function () {
      });
    }
    callback(err);
  });
};

Dc.prototype.getClient = function (clientName, callback) {
  var self = this;
  this.sensu.getClient(clientName, function (err, result) {
    var client = (err) ? {} : result;
    if (!err) {
      self.sensu.sortHistory(client.history, 'check', 'last_status', function () {
      });
      self.sensu.getTimestamp(client.history, 'last_execution', 'lastCheck', function () {
      });
    }
    callback(err, client);
  });
};

Dc.prototype.getInfo = function (callback) {
  var self = this;
  this.sensu.get('info', function (err, result) {
    self.info = (err) ? {} : result;
    if(!err) {
      self.sensu.version = result.sensu.version;
    }
    callback(err);
  });
};

Dc.prototype.pull = function (next, errorCallback) {
  var self = this;
  async.waterfall([
    self.getInfo.bind(this),
    self.getStashes.bind(this),
    self.getChecks.bind(this),
    self.getClients.bind(this),
    self.getEvents.bind(this),
    function (callback) {
      self.sensu.sortEvents(self.sensu.events, 'name', 'status', callback);
    },
    function (callback) {
      self.sensu.sortClients(self.sensu.clients, self.sensu.events, callback);
    },
    function (callback) {
      self.sensu.sortByKey(self.sensu.checks, 'name', callback);
    },
    function (callback) {
      self.sensu.buildClients(callback);
    }
  ], function (err) {
    if (err) {
      self.health = err.message;
      console.log(moment().format() + ' [error] Sensu API ' + self.name + ' returned "' + err + '"');
      errorCallback(JSON.stringify({
        'type': 'error',
        'content': '<strong>Error!</strong> Sensu API ' + self.name + ' returned "' + err + '"'
      }));
      next();
    }
    else {
      self.health = 'ok';
      next();
    }
  });
};

exports.Dc = Dc;