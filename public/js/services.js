'use strict';

var serviceModule = angular.module('uchiwa.services', []);

/**
* Uchiwa
*/
serviceModule.service('uchiwaBackend', ['$http',
  function($http){
    this.createStash = function (payload) {
      return $http.post('/post_stash', payload);
    };
    this.deleteClient = function (client, dc) {
      return $http.get('/delete_client?id=' + client + '&dc=' + dc );
    };
    this.deleteStash = function (payload) {
      return $http.post('/delete_stash', payload);
    };
    this.getClient = function (client, dc) {
      return $http.get('/get_client?id=' + client + '&dc=' + dc );
    };
    this.getConfig = function () {
      return $http.get('/get_config');
    };
    this.getHealth = function () {
      return $http.get('/health/sensu');
    };
    this.getSensu = function () {
      return $http.get('/get_sensu');
    };
    this.resolveEvent = function (payload) {
      return $http.post('/post_event', payload);
    };
  }
]);

/**
* Clients
*/
serviceModule.service('clientsService', ['$location', 'notification', 'uchiwaBackend', function ($location, notification, uchiwaBackend) {
  this.getCheck = function (id, history) {
    return history.filter(function (item) {
      return item.check === id;
    })[0];
  };
  this.getEvent = function (client, check, events) {
    if (!client || !check || events.constructor.toString().indexOf('Array') === -1) { return null; }
    return events.filter(function (item) {
      return (item.client.name === client && item.check.name === check);
    })[0];
  };
  this.resolveEvent = function (dc, client, check) {
    if (!angular.isObject(client) || !angular.isObject(check)) {
      notification('error', 'Could not resolve this event. Try to refresh the page.');
      console.error('Received:\nclient='+ JSON.stringify(client) + '\ncheck=' + JSON.stringify(check));
      return false;
    }

    var checkName = check.name || check.check;
    var payload = {dc: dc, payload: {client: client.name, check: checkName}};

    uchiwaBackend.resolveEvent(payload)
      .success(function () {
        notification('success', 'The event has been resolved.');
        $location.url(encodeURI('/client/' + dc + '/' + client.name));
      })
      .error(function (error) {
        notification('error', 'The event was not resolved. ' + error);
      });
  };
  this.deleteClient = function (dc, client) {
    uchiwaBackend.deleteClient(client, dc)
      .success(function () {
        notification('success', 'The client has been deleted.');
        $location.url('/clients');
        return true;
      })
      .error(function (error) {
        notification('error', 'Could not delete the client '+ client +'. Is Sensu API running on '+ dc +'?');
        console.error(error);
      });
  };
}]);


/**
* Navbar
*/
serviceModule.service('navbarServices', ['$rootScope', function ($rootScope) {
  // Badges count
  this.countStatuses = function (collection, getStatusCode) {
    var criticals = 0;
    var warnings = 0;
    var unknowns = 0;
    var total = collection.length;

    criticals += collection.filter(function (item) {
      return getStatusCode(item) === 2;
    }).length;
    warnings += collection.filter(function (item) {
      return getStatusCode(item) === 1;
    }).length;
    unknowns += collection.filter(function (item) {
      return getStatusCode(item) > 2;
    }).length;

    collection.warning = warnings;
    collection.critical = criticals;
    collection.total = criticals + warnings;
    collection.unknown = unknowns;
    collection.total = total;
    collection.style = collection.critical > 0 ? 'critical' : collection.warning > 0 ? 'warning' : collection.unknown > 0 ? 'unknown' : 'success';
  };
  this.health = function () {
    var alerts = [];
    angular.forEach($rootScope.health, function(value, key) {
      if (value.output !== 'ok') {
        alerts.push('Datacenter <strong>' + key + '</strong> returned: <em>' + value.output + '</em>');
      }
    });
    $rootScope.alerts = alerts;
  };
}]);

/**
* Routing
*/
serviceModule.service('routingService', ['$location', function ($location) {
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

/**
* Stashes
*/
serviceModule.service('stashesService', ['$rootScope', '$modal', 'notification', 'uchiwaBackend', function ($rootScope, $modal, notification, uchiwaBackend) {
  this.construct = function(item) {
    var check;
    var client;
    var path = [];

    if (angular.isObject(item) && angular.isDefined(item.client) && angular.isDefined(item.check)) { // event
      if (!angular.isObject(item.check)) {
        check = item.check;
      }
      else {
        check = {check: item.check};
      }
      if (angular.isObject(item.check)) {
        client = item.client;
      }
      else {
        client = {name: item.client};
      }
    }
    else if (angular.isObject(item) && angular.isDefined(item.name)) { // client
      client = item;
      check = null;
    }
    else { // unknown
      notification('error', 'Cannot handle this stash. Try to refresh the page.');
      console.error('Cannot handle this stash. Received:\nitem: '+ JSON.stringify(item));
      return false;
    }

    path.push(client.name);

    var checkName = '';
    if (check) {
      if (angular.isObject(check.check)) {
        checkName = check.check.name;
      } else {
        checkName = check;
      }
    }
    path.push(checkName);

    return path;
  };
  this.stash = function (e, i) {
    var items = _.isArray(i) ? i : new Array(i);
    var event = e || window.event;
    event.stopPropagation();

    if (items.length === 0) {
      notification('error', 'No items selected');
    } else {
      var modalInstance = $modal.open({ // jshint ignore:line
        templateUrl: 'partials/stash-modal.html',
        controller: 'StashModalCtrl',
        resolve: {
          items: function () {
            return items;
          }
        }
      });
    }
  };
  this.submit = function (element, item) {
    var isAcknowledged = element.acknowledged;
    var path = this.construct(element);
    if (path[1] !== '') {
      path[1] = '/' + path[1];
    }
    if (angular.isUndefined(item.reason)) {
      item.reason = '';
    }
    path = 'silence/' + path[0] + path[1];
    var data = {dc: element.dc, payload: {}};

    $rootScope.skipRefresh = true;
    if (isAcknowledged) {
      data.payload = {path: path};
      uchiwaBackend.deleteStash(data)
        .success(function () {
          notification('success', 'The stash has been deleted.');
          element.acknowledged = !element.acknowledged;
          return true;
        })
        .error(function (error) {
          notification('error', 'The stash was not created. ' + error);
          console.error(error);
          return false;
        });
    }
    else {
      data.payload = {path: path, content: {'reason': item.reason, 'source': 'uchiwa'}};
      if (item.expiration && item.expiration !== -1){
        data.payload.expire = item.expiration;
      }
      data.payload.content.timestamp = Math.floor(new Date()/1000);
      uchiwaBackend.createStash(data)
        .success(function () {
          notification('success', 'The stash has been created.');
          element.acknowledged = !element.acknowledged;
          return true;
        })
        .error(function (error) {
          notification('error', 'The stash was not created. ' + error);
          console.error(error);
          return false;
        });
    }
  };
  this.deleteStash = function (stash) {
    $rootScope.skipRefresh = true;
    var data = {dc: stash.dc, payload: {path: stash.path}};
    uchiwaBackend.deleteStash(data)
      .success(function () {
        notification('success', 'The stash has been deleted.');
        for (var i=0; $rootScope.stashes; i++) {
          if ($rootScope.stashes[i].path === stash.path) {
            $rootScope.stashes.splice(i, 1);
            break;
          }
        }
        return true;
      })
      .error(function (error) {
        notification('error', 'The stash was not created. ' + error);
        console.error(error);
        return false;
      });
  };
}]);
