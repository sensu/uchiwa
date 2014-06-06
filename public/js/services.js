var serviceModule = angular.module('uchiwa.services', []);

/**
 * Socket.IO
 */
serviceModule.factory('socket', function (socketFactory) {
  var socket = socketFactory();
  socket.forward('sensu');
  socket.forward('client');
  return socket;
});

/**
 * Clients
 */
serviceModule.service('clientsService', function(){
  this.stash = function(e, dcName, client, check){
    var event = e || window.event;
    event.stopPropagation();
    var checkName = (_.isUndefined(check)) ? "" : "/" + check.check;
    var isSilenced = (_.isUndefined(check)) ? client.isSilenced : check.isSilenced;
    var path = "silence/"+ client.name + checkName;
    if(isSilenced){
      var payload = {path: path, content:{}};
      socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
      var icon = "fa-volume-up";
    }
    else {
      var timestamp = Math.floor(Date.now() / 1000);
      var payload = {path: path, content:{"reason": "uchiwa", "timestamp": timestamp}};
      socket.emit('create_stash', JSON.stringify({dc: dcName, payload: payload}));
      var icon = "fa-volume-off";
    }
    if (_.isUndefined(check)){
      client.silenceIcon = icon;
      client.isSilenced = !client.isSilenced;
      return client;
    }
    else {
      check.silenceIcon = icon;
      check.isSilenced = !check.isSilenced;
      return check;
    }
  };
  this.resolve = function(e, dcName, client, check){
    var event = e || window.event;
    event.stopPropagation();
    var payload = {client: client.name, check: check.check};
    socket.emit('resolve_event', JSON.stringify({dc: dcName, payload: payload}));
    check.style = "success";
    check.isActive = "Inactive";
    check.event = false;
    check.output = false;
    check.last_check = "Never";
    return check;
  };
   this.delete = function(dcName, clientName){
    var payload = {path: clientName, content:{}};
    socket.emit('delete_client', JSON.stringify({dc: dcName, payload: payload}));
    return true;
  };
});

/**
 * Events
 */
serviceModule.service('eventsService', function(){
  this.stash = function(e, dcName, currentEvent){
    var event = e || window.event;
    event.stopPropagation();
    var path = "silence/"+ currentEvent.client.name + "/" + currentEvent.check.name;
    if(currentEvent.isSilenced){
      var payload = {path: path, content:{}};
      socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
      var icon = "fa-volume-up";
    }
    else {
      var timestamp = Math.floor(Date.now() / 1000);
      var payload = {path: path, content:{"reason": "uchiwa", "timestamp": timestamp}};
      socket.emit('create_stash', JSON.stringify({dc: dcName, payload: payload}));
      var icon = "fa-volume-off";
    }
    currentEvent.silenceIcon = icon;
    currentEvent.isSilenced = !currentEvent.isSilenced;
    return currentEvent;
  };
});

/**
 * Stashes
 */
serviceModule.service('stashesService', function(){
  this.stash = function(dcName, stash, index){
    var checkName = (_.isNull(stash.check)) ? "" : "/" + stash.check;
    var path = "silence/"+ stash.client + checkName;
    var payload = {path: path, content:{}};
    socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
    return stash;
  };
});