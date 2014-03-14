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
  this.client = {};
  this.checks = {};
}

Sensu.prototype.getRest = function(options, callback){
  this.rest.get(options, function(err, result){
    if(err){
      console.log("Fatal error, Sensu API responded with: " + err);
      callback(true, err);
    }
    else {
      callback(null, result);
    }
  });
}

Sensu.prototype.postRest = function(options, data, callback){
  this.rest.post(options, data, function(err, result){
    if(err){
      console.log("Fatal error, Sensu API responded with: " + err);
      callback(true, err);
    }
    else {
      callback(null, result);
    }
  });
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
  this.getRest(options, callback);
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
  this.getRest(options, callback);
}

Sensu.prototype.getChecks = function(callback){
  var options = {
    host: this.host,
    port: this.port,
    path: '/checks',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  }
  this.getRest(options, callback);
}

Sensu.prototype.getClient = function(name, callback){
  var options = {
    host: this.host,
    port: this.port,
    path: '/clients/'+name+'/history',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  }
  this.getRest(options, callback);
}

Sensu.prototype.postStash = function(data, callback){
  var options = {
    host: this.host,
    port: this.port,
    path: '/stashes',
    method: 'POST',
    timeout: this.timeout,
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  }
  this.postRest(options, data, callback);
}

Sensu.prototype.getTimestamp = function(resultSet, unit, callback){
  var now = new Date().getTime();
  async.each(resultSet, function(item, next){
    var timestamp = new Date(item[unit]*1000);
    var seconds = Math.floor((now - timestamp) / 1000);
    if (seconds < 60){
      item.last_check = "1 minute";
    }
    else if (item[unit] == 0){
      item.last_check = "Never";
    }
    else {
      minutes = Math.floor(seconds / 60);
      if (minutes < 60){
        item.last_check = minutes + " minutes";
      }
      else {
        hours = Math.floor(minutes / 24);
        if (hours < 24){
          item.last_check = hours + " hours";
        }
        else {
           item.last_check = Math.floor(hours / 24) + " days";
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
      }
      else if(a.status < b.status){
        return 1;
      }
      else {
        return 0;
      }
    });
    callback();
  });
}

Sensu.prototype.sortEvents = function (resultSet, key, callback){
  resultSet.sort(function(a,b){
    if(a[key] > b[key]){
      return -1;
    }
    else if(a[key] < b[key]){
      return 1;
    }
    else {
      return 0;
    }
  });
  callback();
}

exports.Sensu = Sensu;