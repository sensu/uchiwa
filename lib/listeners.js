'use strict';

var emitters = require('./emitters.js');
var _ = require('underscore');

var listeners = {};

listeners.getDc = function (data, datacenters, callback) {
  if (datacenters.length === 0) {
    return callback('No datacenters found.');
  }
  var dc = datacenters.filter(function (e) {
    return e.name === data.dc;
  });
  if (dc.length !== 1) {
    return callback('The datacenter ' + data.dc + ' was not found.');
  }
  if (_.has(data, 'client')) {
    if (dc[0].sensu.clients.length === 0) {
      return callback('No clients found.');
    }
    var client = dc[0].sensu.clients.filter(function (e) {
      return e.name === data.client;
    });
    if (client.length !== 1) {
      return callback('The client ' + data.client + ' was not found.');
    }
  }
  callback(null, dc[0]);
};

listeners.listen = function (app, sensu, datacenters, publicConfig) {
  var self = this;

  app.io.route('get_sensu', function (req) {
    emitters.send(req, false, sensu, 'sensu');
  });

  app.io.route('get_client', function (req) {
    self.getDc(req.data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(req, err, 'generic');
      }
      else {
        result.getClient(req.data.client, function (err, result) {
          emitters.send(req, err, result, 'client');
        });
      }
    });
  });

  app.io.route('delete_client', function (req) {
    var data = JSON.parse(req.data);
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(req, err, 'generic');
      }
      else {
        result.sensu.delete('clients', data.payload, function (err) {
          emitters.alert(req, err, 'deleteClient');
        });
      }
    });
  });

  app.io.route('create_stash', function (req) {
    var data = JSON.parse(req.data);

    // Set timestamp
    var timestamp = Math.floor(new Date()/1000);
    data.payload.content.timestamp = timestamp;

    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(req, err, 'generic');
      }
      else {
        result.sensu.post('stashes', JSON.stringify(data.payload), function (err) {
          emitters.alert(req, err, 'createStash');
        });
      }
    });
  });

  app.io.route('delete_stash', function (req) {
    var data = JSON.parse(req.data);
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(req, err, 'generic');
      }
      else {
        result.sensu.delete('stashes', data.payload, function (err) {
          emitters.alert(req, err, 'deleteStash');
        });
      }
    });
  });

  app.io.route('resolve_event', function (req) {
    var data = JSON.parse(req.data);
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(req, err, 'generic');
      }
      else {
        result.sensu.post('resolve', JSON.stringify(data.payload), function (err) {
          emitters.alert(req, err, 'resolveEvent');
        });
      }
    });
  });

  app.io.route('get_info', function (req) {
    emitters.send(req, false, publicConfig, 'info');
  });
};

module.exports = listeners;
