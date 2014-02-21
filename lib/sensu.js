var async = require('async');
var Rest = require('./rest.js').Rest;

function Sensu(config) {
  this.host = config.host;
  this.port = config.port;
  this.path = config.path;
  this.method = 'GET';
  this.timeout = config.timeout;
  this.headers = { 'Content-Type': 'application/json' }
  this.rest = new Rest();
  this.events = {};
  this.clients = {};
}

Sensu.prototype.getClients = function(callback){
  var options = {
    host: this.host,
    port: this.port,
    path: '/clients',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  }
  this.callRest(options, callback);
}

Sensu.prototype.getEvents = function(callback){
  var options = {
    host: this.host,
    port: this.port,
    path: '/events',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  }
  this.callRest(options, callback);
}

Sensu.prototype.callRest = function(options, next){
  this.rest.get(options, function(err, result){
    if(err){
      console.log("Error: Sensu API responded with HTTP code " + err);
      next(true);
    }
    else {
      next(null, result);
    }
  });
}

Sensu.prototype.getTimestamp = function(resultSet, unit, callback){
  var now = new Date().getTime();
  async.each(resultSet, function(item, next){
    var timestamp = new Date(item[unit]*1000);
    var seconds = Math.floor((now - timestamp) / 1000);
    if (seconds < 60){
      item.lastCheck = "< 1 minute";
    }
    else {
      minutes = Math.floor(seconds / 60);
      if (minutes < 60){
        item.lastCheck = minutes + " minute(s)";
      }
      else {
        hours = Math.floor(minutes / 24);
        if (hours < 24){
          item.lastCheck = hours + " hour(s)";
        }
        else {
           item.lastCheck = Math.floor(hours / 24) + " day(s)";
        }
      }
    }
    next();
  }, function(err){
    callback();
  });
}

Sensu.prototype.sortClients = function(clients, events, callback){
  async.each(clients, function(client, next){
    var hasEvent = events.filter(function(event){ return event.client == client.name });
    if (hasEvent.length > 0){
      client.events = hasEvent;
      var isCritical = hasEvent.filter(function(event){ return event.status == "2" });
      if (isCritical.length > 0){
        client.status = 2;
      }
      else {
        client.status = 1;
      }
    }
    else {
      client.status = 0;
    }
    next();
  }, function(err){
    clients.sort(function(a,b){
      if(a.status > b.status){
        return -1;
        console.log(">");
      }
      else if(a.status < b.status){
        return 1;
        console.log("<");
      }
      else {
        return 0;
      }
    });
    callback();
  });
}

exports.Sensu = Sensu;