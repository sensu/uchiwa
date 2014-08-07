'use strict';

var serviceModule = angular.module('uchiwa.services', []);

/**
 * Socket.IO
 */
serviceModule.factory('socket', function (socketFactory) {
  var socket = socketFactory();
  socket.forward('sensu');
  socket.forward('client');
  socket.forward('info');
  return socket;
});

/**
 * Page title
 */
serviceModule.factory('Page', function() {
  var title = 'Uchiwa';
  return {
    title: function() { return title + ' | Uchiwa'; },
    setTitle: function(newTitle) { title = newTitle; }
  };
});

/**
 * Notifications
 */
serviceModule.provider('notification', function () {
  this.setOptions = function (options) {
    if (angular.isObject(options)) {
      window.toastr.options = options;
    }
  };
  this.setOptions({});
  this.$get = function () {
    return function (type, message) {
      window.toastr[type](message);
    };
  };
});

/**
 * Underscore.js
 */
serviceModule.factory('underscore', function () {
  if (angular.isUndefined(window._)) {
    console.log('underscore.js is required');
  } else {
    return window._;
  }
});

/**
 * Clients
 */
serviceModule.service('clientsService', ['socket', '$location', function (socket, $location) {
  this.stash = function (dcName, client, check, expire) {
    var checkName = (angular.isUndefined(check)) ? '' : '/' + check.check;
    var isSilenced = (angular.isUndefined(check)) ? client.isSilenced : check.isSilenced;
    var path = 'silence/' + client.name + checkName;
    var payload;
    var icon;
    if (isSilenced) {
      payload = {path: path, content: {}};
      socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
      icon = 'fa-volume-up';
    }
    else {
      payload = {path: path, content: {'reason': 'uchiwa'}};
      if(expire !== -1){
        payload.expire = expire;
      }
      socket.emit('create_stash', JSON.stringify({dc: dcName, payload: payload}));
      icon = 'fa-volume-off';
    }
    if (angular.isUndefined(check)) {
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
  this.resolve = function (dcName, client, check) {
    var payload = {client: client.name, check: check.check};
    console.log(payload);
    socket.emit('resolve_event', JSON.stringify({dc: dcName, payload: payload}));
    check.style = 'success';
    check.isActive = 'Inactive';
    check.event = false;
    check.output = false;
    check.lastCheck = 'Never';
    return check;
  };
  this.remove = function (dcName, clientName) {
    var payload = {path: clientName, content: {}};
    socket.emit('delete_client', JSON.stringify({dc: dcName, payload: payload}));
    $location.url('/clients');
    return true;
  };
}]);

/**
 * Events
 */
serviceModule.service('eventsService', ['socket', function (socket) {
  this.stash = function (dcName, currentEvent, expire) {
    var path = 'silence/' + currentEvent.client.name + '/' + currentEvent.check.name;
    var payload;
    var icon;
    if (currentEvent.isSilenced) {
      payload = {path: path, content: {}};
      socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
      icon = 'fa-volume-up';
    }
    else {
      payload = {path: path, content: {'reason': 'uchiwa'}};
      if(expire !== -1){
        payload.expire = expire;
      }
      socket.emit('create_stash', JSON.stringify({dc: dcName, payload: payload}));
      icon = 'fa-volume-off';
    }
    currentEvent.silenceIcon = icon;
    currentEvent.isSilenced = !currentEvent.isSilenced;
    return currentEvent;
  };
}]);

/**
 * Stashes
 */
serviceModule.service('stashesService', ['socket', function (socket) {
  this.stash = function (dcName, stash) {
    var payload = {path: stash.path, content: {}};
    socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
    return stash;
  };
}]);

/**
 * Routing
 */
serviceModule.service('routingService', ['socket', '$location', function (path, $location) {
  this.go = function (path) {
    path = encodeURI(path);
    $location.url(path);
  };
  this.search = function (e, path) {
    var event = e || window.event;
    event.stopPropagation();
    path = encodeURI(path);
    $location.search(path);
  };
}]);