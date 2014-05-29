var async = require('async');
var moment = require('moment');
var Rest = require('./rest.js').Rest;
var _ = require('underscore');

function Sensu(config) {
  this.host = config.host;
  this.ssl = config.ssl;
  this.port = config.port;
  this.path = config.path;
  this.method = 'GET';
  this.timeout = config.timeout;
  this.headers = { 'Content-Type': 'application/json' };
  this.rest = new Rest();
  this.events = {};
  this.clients = {};
  this.client = {};
  this.checks = {};
  this.stashes = {};
  this.stats = {};
  this.config = config;
}

Sensu.prototype.getRest = function(options, callback){
  this.rest.get(options, this.config, function(err, result){
    if(err){
      console.log("Fatal error, Sensu API responded with: " + err);
      callback("Fatal error while talking with the Sensu API!");
    }
    else {
      callback(null, result);
    }
  });
};

Sensu.prototype.postRest = function(options, data, callback){
  this.rest.post(options, data, this.config, function(err, result){
    if(err){
      console.log("Fatal error, Sensu API responded with: " + err);
      callback(err);
    }
    else {
      callback(null, result);
    }
  });
};

Sensu.prototype.deleteRest = function(options, callback){
  this.rest.delete(options, this.config, function(err){
    if(err){
      console.log("Fatal error, Sensu API responded with: " + err);
      callback(err);
    }
    else {
      callback(null);
    }
  });
};

Sensu.prototype.getClient = function(name, callback){
  var self = this;
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/clients/'+name+'/history',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  // Get client history
  this.getRest(options, function(err, result){
    self.buildClient(name, result, callback);
  });
};

Sensu.prototype.getClients = function(callback){
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/clients',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  this.getRest(options, callback);
};

Sensu.prototype.getEvents = function(callback){
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/events',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  this.getRest(options, callback);
};

Sensu.prototype.getChecks = function(callback){
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/checks',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  this.getRest(options, callback);
};

Sensu.prototype.getStashes = function(callback){
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/stashes',
    method: 'GET',
    timeout: this.timeout,
    headers: this.headers
  };
  this.getRest(options, callback);
};

Sensu.prototype.postStash = function(data, callback){
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/stashes',
    method: 'POST',
    timeout: this.timeout,
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  };
  this.postRest(options, data, callback);
};

Sensu.prototype.deleteStash = function(data, callback){
  var obj = JSON.parse(data)
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/stashes/'+obj.path,
    method: 'DELETE',
    timeout: this.timeout,
    headers: this.headers
  };
  this.deleteRest(options, callback);
};

Sensu.prototype.resolveEvent = function(data, callback){
  var obj = JSON.parse(data)
  var options = {
    host: this.host,
    ssl: this.ssl,
    port: this.port,
    path: this.path+'/resolve',
    method: 'POST',
    timeout: this.timeout,
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  };
  this.postRest(options, data, callback);
};

Sensu.prototype.getTimestamp = function(resultSet, unit, attribute, callback){
  async.each(resultSet, function(item, next){
    if(item[unit] == 0){
      item[attribute] = "Never";
    }
    else {
      var timestamp = (_.isUndefined(item.content)) ? new Date(item[unit]*1000) : new Date(item.content[unit]*1000);
      item[attribute] = moment(timestamp).format('YYYY[-]MM[-]DD HH[:]mm');
    }
    next();
  }, function(err){
    callback();
  });
};

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
      if(a.status > b.status) return -1;
      if(a.status < b.status) return 1;
      return (a.name.toUpperCase() > b.name.toUpperCase()) ? 1 : (a.name.toUpperCase() < b.name.toUpperCase()) ? -1 : 0;
    });
    callback();
  });
};

Sensu.prototype.sortEvents = function (resultSet, id, status, callback){
  resultSet.sort(function(a,b){
    if(a[status] > b[status]) return -1;
    if(a[status] < b[status]) return 1;
    return (a[id].toUpperCase() > b[id].toUpperCase()) ? 1 : (a[id].toUpperCase() < b[id].toUpperCase()) ? -1 : 0;
  });
  callback();
};

Sensu.prototype.sortByKey = function (resultSet, key, callback){
  resultSet.sort(function(a,b){
    if(a[key] > b[key]) return 1;
    if(a[key] < b[key]) return -1;
    return 0;
  });
  callback();
};

Sensu.prototype.buildChecks = function(callback){
  async.each(this.checks, function(check, next){
    check.hasSubscribers = (_.isUndefined(check.subscribers)) ? false : true;
    check.standalone = (_.isUndefined(check.standalone)) ? false : true;
    next();
  },
  function(err){
    callback();    
  });
};

Sensu.prototype.buildClients = function(callback){
  var self = this;
  var isSilenced = function(name){
    if (self.stashes.length == 0) return false;
    var path = "silence/" + name;
    var result = self.stashes.filter(function (e){ return e.path === path });
    return (result.length > 0) ? true : false;
  };

  async.each(this.clients, function(client, next){
    client.style = (client.status == 0) ? "success" : (client.status == 1) ? "warning" : "danger";
    client.eventsSummary = (typeof client.events === 'undefined') ? "" : (client.events.length != 1) ? client.events[0].check + " and " + (client.events.length - 1) + " more..." : client.events[0].check;
    client.isSilenced = isSilenced(client.name);
    client.silenceIcon = (isSilenced(client.name)) ? "fa-volume-off" : "fa-volume-up";
    next();
  },
  function(err){
    callback();    
  });
};

Sensu.prototype.buildClient = function(name, history, callback){
  var self = this;
  var isSilenced = function(name, check){
    if (self.stashes.length == 0) return false;
    var path = "silence/" + name + "/" + check;
    var result = self.stashes.filter(function (e){ return e.path === path });
    return (result.length > 0) ? true : false;
  };
  var hasEvent = function(name, check){
    if(check.last_status == 0 || _.isUndefined(client.events)) return false;
    var event = client.events.filter(function (e) { return e.check == check });
    return (event.length == 0) ? false : event[0];
  };
  var getCheck = function(check){
    if (check == "keepalive") return {name: check, command: "keepalive", subscribers: "", interval: ""};
    var model = self.checks.filter(function (e) { return e.name == check });
    return (model.length == 0) ? {name: check, command: "", subscribers: "", interval: ""} : model[0];
  };

  var find = this.clients.filter(function (e) { return e.name == name });
  if (find.length == 0) callback("No client found with the name " + name, null);
  var client = find[0];
  client.history = history;

  async.each(client.history, function(check, next){
    check.style = (check.last_status == 0) ? "success" : (check.last_status == 1) ? "warning" : "danger";
    check.isSilenced = isSilenced(name, check.check);
    check.silenceIcon = (check.isSilenced) ? "fa-volume-off" : "fa-volume-up";
    check.isActive = (check.last_execution) ? "Active" : "Inactive";
    check.event = hasEvent(client, check.check);
    check.output = (check.event) ? check.event.output : false;
    check.model = getCheck(check.check);
    next();
  },
  function(err){
    callback(null, client);    
  });
};

Sensu.prototype.buildEvents = function(callback){
  var self = this;
  var isSilenced = function(name, check){
    if (self.stashes.length == 0) return false;
    var path = "silence/" + name + "/" + check;
    var result = self.stashes.filter(function (e){ return e.path === path });
    return (result.length > 0) ? true : false;
  };
  async.each(this.events, function(event, next){
    event.style = (event.status == 0) ? "success" : (event.status == 1) ? "warning" : "danger";
    event.isSilenced = isSilenced(event.client, event.check);
    event.silenceIcon = (event.isSilenced) ? "fa-volume-off" : "fa-volume-up";
    next();
  },
  function(err){
    callback();    
  });
};

Sensu.prototype.buildStashes = function(callback){
  async.each(this.stashes, function(stash, next){
    var path = stash.path.split('/');
    stash.client = (_.isUndefined(path[1])) ? null : path[1];
    stash.check = (_.isUndefined(path[2])) ? null : path[2];
    next();
  },
  function(err){
    callback();    
  });
};

exports.Sensu = Sensu;