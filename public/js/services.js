var serviceModule = angular.module('uchiwa.services', []);

/**
 * Socket.IO
 */
serviceModule.factory('socket', function (socketFactory) {
  var socket = socketFactory();
  socket.forward('sensu');
  socket.forward('client');
  socket.forward('stats');
  return socket;
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
 * Utility/Helpers
 */
serviceModule.service('utilityService', ['underscore', function (underscore) {
  // Divide an array into 'n' arrays
  var splitArray = function (array, n) {
    var arrays = [];
    var i, j, temparray, chunk = n;
    for (i = 0, j = array.length; i < j; i += chunk) {
      temparray = array.slice(i, i + chunk);
      arrays.push(temparray);
    }
    return arrays;
  };

  this.getRows = function (array, n) {
    underscore.each(array, function (element, index, list) {
      list[index] = splitArray(element, n);
    });
    return array;
  };
}]);

/**
 * Toggle
 */
serviceModule.service('toggleService', function () {
  var toggle = [];
  this.toggle = toggle;
  this.toggleOn = function (index) {
    if (angular.isUndefined(toggle[index])) {
      toggle[index] = {hidden: false};
    }
    toggle[index].hidden = !toggle[index].hidden;
  };
  this.showOnly = function (index, dc) {
    angular.forEach(dc, function (datacenter, i) {
      if (i === index) {
        toggle[index] = {hidden: false};
      }
      else {
        toggle[i] = {hidden: true};
      }
    });
  };
  this.showAll = function (dc) {
    angular.forEach(dc, function (datacenter, i) {
      toggle[i] = {hidden: false};
    });
  };
});

/**
 * Toggle Client
 */
serviceModule.service('toggleClientService', function () {
  var toggle = [];
  this.toggle = toggle;
  this.toggleOn = function (index) {
    if (angular.isUndefined(toggle[index])) {
      toggle[index] = {hidden: false};
    }
    toggle[index].hidden = !toggle[index].hidden;
  };
});

/**
 * Clients
 */
serviceModule.service('clientsService', ['socket', function (socket) {
  this.stash = function (e, dcName, client, check) {
    var event = e || window.event;
    event.stopPropagation();
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
      var timestamp = Math.floor(Date.now() / 1000);
      payload = {path: path, content: {'reason': 'uchiwa', 'timestamp': timestamp}};
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
  this.resolve = function (e, dcName, client, check) {
    var event = e || window.event;
    event.stopPropagation();
    var payload = {client: client.name, check: check.check};
    socket.emit('resolve_event', JSON.stringify({dc: dcName, payload: payload}));
    check.style = 'success';
    check.isActive = 'Inactive';
    check.event = false;
    check.output = false;
    check.lastCheck = 'Never';
    return check;
  };
  this.delete = function (dcName, clientName) {
    var payload = {path: clientName, content: {}};
    socket.emit('delete_client', JSON.stringify({dc: dcName, payload: payload}));
    return true;
  };
}]);

/**
 * Events
 */
serviceModule.service('eventsService', ['socket', function (socket) {
  this.stash = function (e, dcName, currentEvent) {
    var event = e || window.event;
    event.stopPropagation();
    var path = 'silence/' + currentEvent.client.name + '/' + currentEvent.check.name;
    var payload;
    var icon;
    if (currentEvent.isSilenced) {
      payload = {path: path, content: {}};
      socket.emit('delete_stash', JSON.stringify({dc: dcName, payload: payload}));
      icon = 'fa-volume-up';
    }
    else {
      var timestamp = Math.floor(Date.now() / 1000);
      payload = {path: path, content: {'reason': 'uchiwa', 'timestamp': timestamp}};
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
