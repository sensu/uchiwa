'use strict';

var controllerModule = angular.module('uchiwa.controllers', []);

/**
* Init
*/
controllerModule.controller('init', ['$rootScope', '$scope', '$interval', 'notification', 'Page', 'uchiwaBackend', 'helpers',
  function ($rootScope, $scope, $interval, notification, Page, uchiwaBackend, helpers) {
    $scope.Page = Page;
    $rootScope.skipRefresh = false;
    $rootScope.alerts = [];
    $rootScope.events = [];
    $rootScope.helpers = helpers;

    uchiwaBackend.getConfig()
      .success(function (data) {
        $interval($rootScope.getSensu, data.Uchiwa.Refresh * 1000);
      })
      .error(function () {
        $interval($rootScope.getSensu, 10000);
      });

    $rootScope.getSensu = function () {
      if ($rootScope.skipRefresh) {
        $rootScope.skipRefresh = false;
        return;
      }
      uchiwaBackend.getHealth()
        .success(function (data) {
          $rootScope.health = data;
        });
      uchiwaBackend.getSensu()
        .success(function (data) {
          angular.forEach(data, function(value, key) { // initialize null elements
            if (!value || value === null) {
              data[key] = [];
            }
          });
          $rootScope.checks = data.Checks;
          $rootScope.dc = data.Dc;

          $rootScope.clients = _.map(data.Clients, function(client) {
            var existingClient = _.findWhere($rootScope.clients, {name: client.name, dc: client.dc});
            if (existingClient !== undefined) {
              client = angular.extend(existingClient, client);
            }
            return existingClient || client;
          });

          $rootScope.events = _.map(data.Events, function(event) {
            event._id = event.dc + '/' + event.client.name + '/' + event.check.name;
            var existingEvent = _.findWhere($rootScope.events, {_id: event._id});
            if (existingEvent !== undefined) {
              event = angular.extend(existingEvent, event);
            }
            return existingEvent || event;
          });

          $rootScope.stashes = data.Stashes;
          $rootScope.subscriptions = data.Subscriptions;
          $scope.$broadcast('sensu');

        })
        .error(function (error) {
          notification('error', 'Could not fetch Sensu data. Is Uchiwa running?');
          console.error('Error: '+ JSON.stringify(error));
        });
    };

    $scope.$on('$routeChangeSuccess', function () {
      $rootScope.getSensu();
    });

    $scope.$on('notification', function (type, message) {
      notification(type, message);
    });
  }
]);

/**
* Checks
*/
controllerModule.controller('checks', ['$scope', '$routeParams', 'routingService', 'Page',
  function ($scope, $routeParams, routingService, Page) {
    Page.setTitle('Checks');
    $scope.pageHeaderText = 'Checks';
    $scope.predicate = 'name';

    // Helpers
    $scope.subscribersSummary = function(subscribers){
      return subscribers.join(' ');
    };

    // Routing
    $scope.filters = {};
    routingService.initFilters($routeParams, $scope.filters, ['dc', 'limit', 'q']);
    $scope.$on('$locationChangeSuccess', function(){
      routingService.updateFilters($routeParams, $scope.filters);
    });

    // Services
    $scope.permalink = routingService.permalink;

  }
]);

/**
* Client
*/
controllerModule.controller('client', ['$scope', '$routeParams', 'clientsService', 'notification', 'Page', 'routingService', 'stashesService', 'uchiwaBackend',
  function ($scope, $routeParams, clientsService, notification, Page, routingService, stashesService, uchiwaBackend) {

    $scope.predicate = '-last_status';
    $scope.missingClient = false;

    // Retrieve client
    $scope.clientId = decodeURI($routeParams.clientId);
    $scope.dcId = decodeURI($routeParams.dcId);
    $scope.pull = function() {
      uchiwaBackend.getClient($scope.clientId, $scope.dcId)
        .success(function (data) {
          $scope.$emit('client', data);
        })
        .error(function (error) {
          // Stop the pulling interval and set scope to display an error message
          clearTimeout(timer);
          $scope.missingClient = true;
          console.error('Error: '+ JSON.stringify(error));
        });
    };

    $scope.pull();
    var timer = setInterval($scope.pull, 10000);

    $scope.$on('client', function (event, data) {
      $scope.client = data;
      $scope.pageHeaderText = $scope.client.name;

      // Retrieve check & event
      $scope.requestedCheck = decodeURI($routeParams.check);
      $scope.selectedCheck = getCheck($scope.requestedCheck, $scope.client.history);
      $scope.selectedEvent = getEvent($scope.client.name, $scope.requestedCheck, $scope.events);

      // Set page title
      if(angular.isDefined($scope.selectedCheck)) {
        Page.setTitle($scope.requestedCheck + ' - ' + $scope.client.name);
      }
      else {
        Page.setTitle($scope.client.name);
      }
    });

    // Routing
    $scope.$on('$routeUpdate', function(){
      // Retrieve check & event
      $scope.requestedCheck = decodeURI($routeParams.check);
      $scope.selectedCheck = getCheck($scope.requestedCheck, $scope.client.history);
      $scope.selectedEvent = getEvent($scope.client.name, $scope.requestedCheck, $scope.events);

      if(angular.isDefined($scope.selectedCheck)) {
        Page.setTitle($scope.requestedCheck + ' - ' + $scope.client.name);
      }
      else {
        Page.setTitle($scope.client.name);
      }
    });

    $scope.$on('$destroy', function() {
      clearInterval(timer);
    });

    // Sanitize - only display useful information 'acknowledged', 'dc', 'events', 'eventsSummary', 'history', 'status', 'timestamp'
    /* jshint ignore:start */
    var clientWhitelist = [ 'acknowledged', 'dc', 'events', 'eventsSummary', 'history', 'output', 'status', 'timestamp' ];
    var checkWhitelist = [ 'dc', 'hasSubscribers', 'name'];
    $scope.sanitizeObject = function(type, key){
      return eval(type + 'Whitelist').indexOf(key) === -1;
    };
    /* jshint ignore:end */

    // Services
    $scope.deleteClient = clientsService.deleteClient;
    $scope.resolveEvent = clientsService.resolveEvent;
    $scope.permalink = routingService.permalink;
    $scope.stash = stashesService.stash;
    var getCheck = clientsService.getCheck;
    var getEvent = clientsService.getEvent;
  }
]);

/**
* Clients
*/
controllerModule.controller('clients', ['$scope', '$rootScope', '$routeParams', 'routingService', 'stashesService', 'clientsService', '$filter', 'Page',
  function ($scope, $rootScope, $routeParams, routingService, stashesService, clientsService, $filter, Page) {
    Page.setTitle('Clients');
    $scope.pageHeaderText = 'Clients';
    $scope.predicate = '-status';

    // Routing
    $scope.filters = {};
    routingService.initFilters($routeParams, $scope.filters, ['dc', 'subscription', 'limit', 'q']);
    $scope.$on('$locationChangeSuccess', function(){
      routingService.updateFilters($routeParams, $scope.filters);
    });

    // Services
    $scope.go = routingService.go;
    $scope.permalink = routingService.permalink;
    $scope.stash = stashesService.stash;
    $scope.deleteClient = clientsService.deleteClient;

    // Helpers
    $scope.selectedClients = function(clients) {
      return _.filter(clients, function(client) {
        return client.selected === true;
      });
    };

    $scope.selectClients = function(selectModel) {
      var filteredClients = $filter('filter')($rootScope.clients, $scope.filters.q);
      filteredClients = $filter('filter')(filteredClients, {dc: $scope.filters.dc});
      filteredClients = $filter('hideSilenced')(filteredClients, $scope.filters.silenced);
      _.each(filteredClients, function(client) {
        client.selected = selectModel.selected;
      });
    };

    $scope.deleteClients = function(clients) {
      _.each(clients, function(client) {
        $scope.deleteClient(client.dc, client.name);
      });
    };

    $scope.$watch('filters.q', function(newVal) {
      var matched = $filter('filter')($rootScope.clients, '!'+newVal);
      _.each(matched, function(match) {
        match.selected = false;
      });
    });

    $scope.$watch('filters.dc', function(newVal) {
      var matched = $filter('filter')($rootScope.clients, {dc: '!'+newVal});
      _.each(matched, function(match) {
        match.selected = false;
      });
    });

    $scope.$watch('filters.silenced', function() {
      var matched = $filter('filter')($rootScope.clients, {acknowledged: true});
      _.each(matched, function(match) {
        match.selected = false;
      });
    });
  }
]);

/**
* Events
*/
controllerModule.controller('events', ['$cookieStore', '$scope', '$rootScope', '$routeParams','routingService', 'settings', 'stashesService', 'clientsService', '$filter', 'Page',
  function ($cookieStore, $scope, $rootScope, $routeParams, routingService, settings, stashesService, clientsService, $filter, Page) {
    Page.setTitle('Events');
    $scope.pageHeaderText = 'Events';
    $scope.predicate = '-check.status';
    $scope.filters = {};

    // Routing
    routingService.initFilters($routeParams, $scope.filters, ['dc', 'limit', 'q']);
    $scope.$on('$locationChangeSuccess', function(){
      routingService.updateFilters($routeParams, $scope.filters);
    });

    // Services
    $scope.go = routingService.go;
    $scope.permalink = routingService.permalink;
    $scope.stash = stashesService.stash;
    $scope.resolveEvent = clientsService.resolveEvent;

    // Hide silenced
    $scope.filters.silenced = $cookieStore.get('hideSilenced') || settings.hideSilenced;
    $scope.$watch('filters.silenced', function () {
      $cookieStore.put('hideSilenced', $scope.filters.silenced);
    });

    // Helpers
    $scope.selectedEvents = function(events) {
      return _.filter(events, function(event) {
        return event.selected === true;
      });
    };

    $scope.selectEvents = function(selectModel) {
      var filteredEvents = $filter('filter')($rootScope.events, $scope.filters.q);
      filteredEvents = $filter('filter')(filteredEvents, {dc: $scope.filters.dc});
      filteredEvents = $filter('hideSilenced')(filteredEvents, $scope.filters.silenced);
      _.each(filteredEvents, function(event) {
        event.selected = selectModel.selected;
      });
    };

    $scope.resolveEvents = function(events) {
      _.each(events, function(event) {
        $scope.resolveEvent(event.dc, event.client, event.check);
      });
    };

    $scope.$watch('filters.q', function(newVal) {
      var matched = $filter('filter')($rootScope.events, '!'+newVal);
      _.each(matched, function(match) {
        match.selected = false;
      });
    });

    $scope.$watch('filters.dc', function(newVal) {
      var matched = $filter('filter')($rootScope.events, {dc: '!'+newVal});
      _.each(matched, function(match) {
        match.selected = false;
      });
    });

    $scope.$watch('filters.silenced', function() {
      var matched = $filter('filter')($rootScope.events, {acknowledged: true});
      _.each(matched, function(match) {
        match.selected = false;
      });
    });
  }
]);

/**
* Info
*/
controllerModule.controller('info', ['$scope', 'notification', 'Page', 'version', 'uchiwaBackend',
  function ($scope, notification, Page, version, uchiwaBackend) {
    $scope.pageHeaderText = 'Info';
    $scope.uchiwa = {};
    Page.setTitle('Info');
    $scope.uchiwa.version = version.uchiwa;

    uchiwaBackend.getConfig()
      .success(function (data) {
        $scope.uchiwa.config = JSON.stringify(data, null, 2);
      })
      .error(function (error) {
        notification('error', 'Could not fetch Uchiwa config. Is Uchiwa running?');
        console.error('Error: '+ JSON.stringify(error));
      });
  }
]);

/**
* Navbar
*/
controllerModule.controller('navbar', ['$rootScope', '$scope', 'navbarServices', 'routingService',
  function ($rootScope, $scope, navbarServices, routingService) {

    // Services
    $scope.go = routingService.go;

    $scope.$on('sensu', function () {
      // Update badges
      navbarServices.countStatuses($rootScope.clients, function (item) {
        return item.status;
      });
      navbarServices.countStatuses($rootScope.events, function (item) {
        return item.check.status;
      });

      // Update alert badge
      navbarServices.health();
    });
  }
]);

/**
* Settings
*/
controllerModule.controller('settings', ['$cookies', '$scope', 'Page',
  function ($cookies, $scope, Page) {
    $scope.pageHeaderText = 'Settings';
    Page.setTitle('Settings');
    $scope.$watch('currentTheme', function (theme) {
      $scope.$emit('theme:changed', theme);
    });
  }
]);

/**
* Sidebar
*/
controllerModule.controller('sidebar', ['$scope', '$location',
  function ($scope, $location) {
    $scope.getClass = function(path) {
      if ($location.path().substr(0, path.length) === path) {
        return 'selected';
      } else {
        return '';
      }
    };
  }
]);

/**
* Stashes
*/
controllerModule.controller('stashes', ['$scope', '$routeParams', 'routingService', 'stashesService', 'Page',
  function ($scope, $routeParams, routingService, stashesService, Page) {
    Page.setTitle('Stashes');
    $scope.pageHeaderText = 'Stashes';
    $scope.predicate = 'client';

    // Helpers
    //$scope.deleteStash = function (dcName, stash, index) {
    //  stashesService.deleteStash(stash);
      //$scope.stashes.splice(index, 1);
    //};
    $scope.deleteStash = stashesService.deleteStash;

    // Routing
    $scope.filters = {};
    routingService.initFilters($routeParams, $scope.filters, ['dc', 'limit', 'q']);
    $scope.$on('$locationChangeSuccess', function(){
      routingService.updateFilters($routeParams, $scope.filters);
    });

    // Services
    $scope.permalink = routingService.permalink;

  }
]);

/**
* Stash Modal
*/
controllerModule.controller('StashModalCtrl', ['$scope', '$filter', '$modalInstance', 'items', 'stashesService',
  function ($scope, $filter, $modalInstance, items, stashesService) {
    $scope.items = items;
    $scope.acknowledged = $filter('filter')(items, {acknowledged: true}).length;
    $scope.itemType = items[0].hasOwnProperty('client') ? 'check' : 'client';
    $scope.stash = {};
    $scope.stash.expirations = {
      '900': 900,
      '3600': 3600,
      '86400': 86400,
      'none': -1
    };
    $scope.stash.reason = '';
    $scope.stash.expiration = 900;

    $scope.stashForItem = function(stashes, item) {
      var path = 'silence/';

      if ($scope.itemType === 'client') {
        path = path + item.name;
      } else if ($scope.itemType === 'check') {
        path = path + item.client.name + '/' + item.check.name;
      }

      return _.findWhere(stashes, {
        dc: item.dc,
        path: path
      });
    };

    $scope.ok = function () {
      _.each(items, function(item) {
        stashesService.submit(item, $scope.stash);
      });
      $modalInstance.close();
    };
    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }
]);
