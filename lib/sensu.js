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

/**
 * Getters
 */
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

/**
 * Finders
 */
Sensu.prototype.findCheck = function(name){
  if (name == "keepalive") return {name: name, command: "keepalive", subscribers: "", interval: ""};
  var check = this.checks.filter(function (e) { return e.name == name });
  return (check.length == 0) ? {name: name, command: "", subscribers: "", interval: ""} : check[0];
};

Sensu.prototype.findClient = function(name){
  if(this.clients.length == 0) return null;
  var client = this.clients.filter(function (e) { return e.name == name });
  return (client.length == 0) ? null : client[0];
};

Sensu.prototype.findEvents = function(client){
  if(this.events.length == 0) return null;
  var eventsFound = this.events.filter(function(event){ return event.client.name == client.name });
  return (eventsFound.length == 0) ? null : eventsFound;
};

Sensu.prototype.findStash = function(clientName, checkName){
  if (this.stashes.length == 0) return false;
  var check = (_.isUndefined(checkName)) ? "" : "/" + checkName;
  var path = "silence/" + clientName + check;
  var result = this.stashes.filter(function (e){ return e.path === path });
  return (result.length > 0) ? true : false;
};

/**
 * Sorters
 */
Sensu.prototype.sortClients = function(clients, events, callback){
  self = this;
  async.each(clients, function(client, next){
    client.events = self.findEvents(client);
    if (client.events && client.events.length > 0){
      var isCritical = client.events.filter(function(event){ return event.check.status == "2" });
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

Sensu.prototype.sortHistory = function (resultSet, id, status, callback){
  resultSet.sort(function(a,b){
    if(a[status] > b[status]) return -1;
    if(a[status] < b[status]) return 1;
    return (a[id].toUpperCase() > b[id].toUpperCase()) ? 1 : (a[id].toUpperCase() < b[id].toUpperCase()) ? -1 : 0;
  });
  callback();
};

Sensu.prototype.sortEvents = function (resultSet, id, status, callback){
  if(resultSet.length == 0 || !_.isArray(resultSet) ) return callback();
  resultSet.sort(function(a,b){
    if(a.check[status] > b.check[status]) return -1;
    if(a.check[status] < b.check[status]) return 1;
    return (a.check[id].toUpperCase() > b.check[id].toUpperCase()) ? 1 : (a.check[id].toUpperCase() < b.check[id].toUpperCase()) ? -1 : 0;
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


/**
 * Builders
 */
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
  async.each(this.clients, function(client, next){
    client.style = (client.status == 0) ? "success" : (client.status == 1) ? "warning" : "danger";
    client.eventsSummary = (!client.events) ? "" : (client.events.length != 1) ? client.events[0].check.name + " and " + (client.events.length - 1) + " more..." : client.events[0].check.name;
    client.isSilenced = self.findStash(client.name);
    client.silenceIcon = (client.isSilenced) ? "fa-volume-off" : "fa-volume-up";
    if(!client.version) client.version = "<= 0.12.6";

    // Cleanup client events once we are done to avoid circular structure
    client.events = null;
    next();
  },
  function(err){
    callback();    
  });
};

Sensu.prototype.buildClient = function(name, history, callback){
  var self = this;

  var hasEvent = function(client, check){;
    if(check.last_status == 0 || !client.events) return false;
    var event = client.events.filter(function (e) { return e.check.name == check.check });
    return (event.length == 0) ? false : event[0];
  };

  var client = this.findClient(name);
  client.history = history;
  client.events = self.findEvents(client);

  async.each(client.history, function(check, next){

    check.style = (check.last_status == 0) ? "success" : (check.last_status == 1) ? "warning" : "danger";
    check.isSilenced = self.findStash(name, check.check);
    check.silenceIcon = (check.isSilenced) ? "fa-volume-off" : "fa-volume-up";
    check.isActive = (check.last_execution) ? "Active" : "Inactive";
    check.event = hasEvent(client, check);

    // Build backward compatible object for Sensu < 0.13.0
    if(check.event) check.event = self.eventData(self, check.event);
    
    check.output = (check.event) ? check.event.check.output : false;
    check.model = self.findCheck(check.check);
    next();
  },
  function(err){
    callback(null, client);    
  });
};

// Build backward compatible object for Sensu < 0.13.0
Sensu.prototype.eventData = function(self, event){
  var getTimestamp = function(epoch){
    if (epoch == 0) return "Never";
    var timestamp = new Date(epoch*1000);
    return moment(timestamp).format('YYYY[-]MM[-]DD HH[:]mm');
  };

  var eventData = (_.has(event, "id")) ? event : {};

  // Check if Sensu < 0.13.0
  if(!_.has(event, "id")){
    eventData.id = null;
    var fullClient = this.findClient(event.client);
    eventData.client = {name: fullClient.name, address: fullClient.address, subscription: fullClient.subscription, version: "<= 0.12.6"};
    eventData.check = this.findCheck(event.check);
    eventData.check.status = event.status;
    eventData.check.issued = event.issued;
    eventData.check.executed = null;
    eventData.check.output = event.output;
    eventData.check.duration = null;
    eventData.check.history = null;
    eventData.occurrences = event.occurrences;
    eventData.action = "create";
    eventData.check.issued = event.issued;
  }

  eventData.style = (eventData.check.status == 0) ? "success" : (eventData.check.status == 1) ? "warning" : "danger";
  eventData.isSilenced = this.findStash(eventData.client.name, eventData.check.name);
  eventData.silenceIcon = (eventData.isSilenced) ? "fa-volume-off" : "fa-volume-up";
  eventData.check.last_issued = getTimestamp(eventData.check.issued);
  if(eventData.check.executed) eventData.check.last_executed = getTimestamp(eventData.check.executed);
  return eventData;
};

Sensu.prototype.buildEvents = function(callback){
  var self = this;
  _.each(this.events, function(element, index, list){
    list[index] = self.eventData(self, element);
  });
  callback();
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