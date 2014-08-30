'use strict';

var emitters = require('./emitters.js');
var _ = require('underscore');

var listeners = {};

var clients = {};

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

listeners.listen = function (socket, sensu, datacenters, publicConfig) {
  var self = this;

  // Keep track of active clients
  clients[socket.id] = socket;

  // Remove client on disconnection
  socket.on('disconnect', function () {
    delete clients[socket.id];
  });

  socket.on('get_sensu', function () {
    emitters.send(clients[socket.id], false, sensu, 'sensu');
  });

  socket.on('get_client', function (data) {
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(clients[socket.id], err, 'generic');
      }
      else {
        result.getClient(data.client, function (err, result) {
          emitters.send(clients[socket.id], err, result, 'client');
        });
      }
    });
  });

  socket.on('delete_client', function (data) {
    data = JSON.parse(data);
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(clients[socket.id], err, 'generic');
      }
      else {
        result.sensu.delete('clients', data.payload, function (err) {
          emitters.alert(clients[socket.id], err, 'deleteClient');
        });
      }
    });
  });

  socket.on('create_stash', function (data) {
    data = JSON.parse(data);
   
    // Set timestamp
    var timestamp = Math.floor(new Date()/1000);
    data.payload.content.timestamp = timestamp;
   
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(clients[socket.id], err, 'generic');
      }
      else {
        result.sensu.post('stashes', JSON.stringify(data.payload), function (err) {
          emitters.alert(clients[socket.id], err, 'createStash');
        });
      }
    });
  });

  socket.on('delete_stash', function (data) {
    data = JSON.parse(data);
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(clients[socket.id], err, 'generic');
      }
      else {
        result.sensu.delete('stashes', data.payload, function (err) {
          emitters.alert(clients[socket.id], err, 'deleteStash');
        });
      }
    });
  });

  socket.on('resolve_event', function (data) {
    data = JSON.parse(data);
    self.getDc(data, datacenters, function (err, result) {
      if (err) {
        emitters.alert(clients[socket.id], err, 'generic');
      }
      else {
        result.sensu.post('resolve', JSON.stringify(data.payload), function (err) {
          emitters.alert(clients[socket.id], err, 'resolveEvent');
        });
      }
    });
  });

  socket.on('get_info', function () {
    emitters.send(clients[socket.id], false, publicConfig, 'info');
  });
};

module.exports = listeners;