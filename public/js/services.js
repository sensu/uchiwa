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
  var filtersDefaultValues = {
    'limit': 50
  };
  this.go = function (path) {
    path = encodeURI(path);
    $location.url(path);
  };
  this.deleteEmptyParameter = function (routeParams, key) {
    if (routeParams[key] === '') {
      delete $location.$$search[key];
      $location.$$compose();
    }
  };
  this.initFilters = function (routeParams, filters, possibleFilters) {
    var self = this;
    angular.forEach(possibleFilters, function (key) {
      if (angular.isDefined(routeParams[key])) {
        self.updateValue(filters, routeParams[key], key);
        self.deleteEmptyParameter(routeParams, key);
      }
      else {
        self.updateValue(filters, '', key);
      }
    });
  };
  this.permalink = function (e, key, value) {
    //var event = e || window.event;
    //event.stopPropagation();
    $location.search(key, value);
  };
  this.updateFilters = function (routeParams, filters) {
    var self = this;
    angular.forEach(routeParams, function (value, key) {
      self.updateValue(filters, value, key);
      self.deleteEmptyParameter(routeParams, key);
    });
  };
  this.updateValue = function (filters, value, key) {
    if (value === '') {
      filters[key] = filtersDefaultValues[key] ? filtersDefaultValues[key] : value;
    }
    else {
      filters[key] = value;
    }
  };
}]);