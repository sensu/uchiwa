'use strict';

var async = require('async');
var logger = require('./logger.js');
var _ = require('underscore');
var Sensu = require('./sensu.js').Sensu;

function Dc(config) {
  this.name = config.name;
  this.sensu = new Sensu(config);
  this.status = 0;
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
  this.status = (this.criticals > 0) ? 2 : (this.warnings > 0) ? 1 : (this.unknown > 0) ? 3 : 0;
};

Dc.prototype.get = function (endpoint, callback) {
  var self = this;
  this.sensu.get(endpoint, function (err, result) {
    self.sensu[endpoint] = (err) ? [] : result;
    callback(err);
  });
};

Dc.prototype.getClient = function (clientName, callback) {
  this.sensu.getClient(clientName, function (err, result) {
    var client = (err) ? {} : result;
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
    self.get.bind(this, 'stashes'),
    self.get.bind(this, 'checks'),
    self.get.bind(this, 'clients'),
    self.get.bind(this, 'events'),
    function (callback) {
      self.sensu.buildEvents(callback);
    },
    function (callback) {
      self.sensu.buildClients(callback);
    }
  ], function (err) {
    if (err) {
      self.health = err.message;
      logger.error('Processing Sensu API ' + self.name + ' returned "' + err + '"');
      errorCallback(JSON.stringify({
        'type': 'error',
        'content': '<strong>Error</strong> with Sensu API ' + self.name
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
