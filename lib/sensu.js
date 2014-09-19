'use strict';

var _ = require('underscore');
var async = require('async');
var Rest = require('./rest.js').Rest;

function Sensu(config) {
  this.host = config.host;
  this.ssl = config.ssl || false;
  this.port = config.port;
  this.path = config.path || '';
  this.method = 'GET';
  this.timeout = config.timeout || 5000;
  this.headers = {'Content-Type': 'application/json'};
  this.rest = new Rest();
  this.events = [];
  this.clients = [];
  this.client = {};
  this.checks = [];
  this.stashes = [];
  this.version = '';
  this.config = config;
}

/**
 * Getters
 */
Sensu.prototype.getClient = function (name, callback) {
  var self = this;
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path + '/clients/' + name + '/history',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  // Get client history
  this.rest.get(options, this.config, function (err, result) {
    if (!err) {
      self.buildClient(name, result, callback);
    }
    else {
      callback(err, result);
    }
  });
};

Sensu.prototype.get = function (endpoint, callback) {
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path + '/' + endpoint,
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  this.rest.get(options, this.config, callback);
};

Sensu.prototype.delete = function (endpoint, data, callback) {
  var obj = data;
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path + '/' + endpoint + '/' + obj.path,
    method: 'DELETE',
    timeout: this.timeout,
    headers: this.headers
  };
  this.rest.delete(options, this.config, callback);
};

Sensu.prototype.post = function (endpoint, data, callback) {
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path + '/' + endpoint,
    method: 'POST',
    timeout: this.timeout,
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  };
  this.rest.post(options, data, this.config, callback);
};

/**
 * Finders
 */
Sensu.prototype.findStash = function (clientName, checkName) {
  if (!_.isArray(this.stashes) || this.stashes.length === 0) {
    return null;
  }
  var check = (_.isUndefined(checkName)) ? '' : '/' + checkName;
  var path = 'silence/' + clientName + check;
  var result = this.stashes.filter(function (e) {
    return e.path === path;
  });
  return (result.length === 0) ? null : true;
};

/**
 * Builders
 */

Sensu.prototype.buildClient = function (name, history, callback) {
  var self = this;

  var getClient = function (name, next) {
    if (!_.isArray(self.clients) || self.clients.length === 0) { return next(true); }
    var client = _.filter(self.clients, function (c) { return c.name === name; });
    return (client.length === 0) ? next(true) : next(null, client[0]);
  };

  var getModel = function (check) {
    if (!_.isArray(self.checks) || self.checks.length === 0) { check.model = null; return check; }
    var model = _.filter(self.checks, function (c) { return c.name === check.check; });
    if (model.length === 0) { check.model = null; return check; }
    check.model = model[0];
    return check;
  };

  var getOutput = function (check) {
    if (check.last_status === 0 || !_.isArray(self.events) || self.events.length === 0) { check.output = ''; return check; } // jshint ignore:line
    var event = _.filter(self.events, function (event) { return (event.client.name === name && event.check.name === check.check); });
    if (event.length === 0) { check.output = ''; return check; }
    check.output = event[0].check.output;
    return check;
  };

  getClient(name, function(err, result) {
    if (err) { callback('Did not found any client named ' + name); }
    var client = result;
    client.history = history;

    async.each(client.history, function (check, next) {
      check.acknowledged = self.findStash(name, check.check);
      getOutput(check);
      getModel(check);
      next();
    },
    function () {
      callback(null, client);
    });
  });
};

Sensu.prototype.buildClients = function (callback) {
  var self = this;

  var getStatus = function(client, callback) {
    if (!_.isArray(self.events) || self.events.length === 0) {
      client.status = 0;
      return callback([]);
    }
    var events = _.filter(self.events, function (event) { return event.client.name === client.name; });
    if (events.length === 0) {
      client.status = 0;
      return callback([]);
    }
    var criticals = _.filter(events, function (event) { return event.check.status === 2; });
    if (criticals.length > 0) {
      client.status = 2;
      return callback(events);
    }
    var warnings = _.filter(events, function (event) { return event.check.status === 1; });
    client.status = (warnings.length > 0) ? 1 : 3;
    return callback(events);
  };

  var findEvents = function (client) {
    getStatus(client, function (result) {
      client.eventsSummary = (result.length === 0) ? '' : (result.length !== 1) ? result[0].check.name + ' and ' + (result.length - 1) + ' more...' : result[0].check.name;
    });
  };

  async.each(this.clients, function (client, next) {
    client.version = client.version || '0.12.x';
    findEvents(client);
    client.acknowledged = self.findStash(client.name);
    next();
  },
  function () {
    return callback();
  });
};

Sensu.prototype.buildEvents = function (callback) {
  var self = this;
  _.each(this.events, function (element, index, list) {
    var event = (_.has(element, 'id')) ? element : {};

    // Build backward compatible event object for Sensu < 0.13.0
    if (!_.has(element, 'id')) {
      event.client = { name: element.client };
      event.check = {};
      event.check.name = element.check;
      var properties = ['issued', 'output', 'status'];
      _.each(properties, function (property) {
        if (element[property]) {
          event.check[property] = element[property];
        }
      });
      event.occurrences = element.occurrences || 1;
      event.action = (element.flapping) ? 'flapping' : 'create';
    }

    event.acknowledged = self.findStash(event.client.name, event.check.name);
    list[index] = event;
  });
  callback();
};

exports.Sensu = Sensu;
