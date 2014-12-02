'use strict';

var controllerModule = angular.module('uchiwa.controllers', []);

/**
* Init
*/
controllerModule.controller('init', ['$rootScope', '$scope', 'notification', 'pollingFactory', 'Page', 'uchiwaBackend',
  function ($rootScope, $scope, notification, pollingFactory, Page, uchiwaBackend) {
    $scope.Page = Page;
    $rootScope.skipRefresh = false;
    $rootScope.alerts = [];

    uchiwaBackend.getConfig()
      .success(function (data) {
        pollingFactory.callFnOnInterval(function () { $rootScope.getSensu(); }, data.Uchiwa.Refresh);
      })
      .error(function () {
        pollingFactory.callFnOnInterval(function () { $rootScope.getSensu(); }, 10);
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
          $rootScope.clients = data.Clients;
          $rootScope.dc = data.Dc;
          $rootScope.events = data.Events;
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
controllerModule.controller('clients', ['$scope', '$routeParams', 'routingService', 'stashesService', 'Page',
  function ($scope, $routeParams, routingService, stashesService, Page) {
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
  }
]);

/**
* Events
*/
controllerModule.controller('events', ['$cookieStore', '$scope', '$routeParams','routingService', 'settings', 'stashesService', 'Page',
  function ($cookieStore, $scope, $routeParams, routingService, settings, stashesService, Page) {
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

    // Hide silenced
    $scope.filters.silenced = $cookieStore.get('hideSilenced') || settings.hideSilenced;
    $scope.$watch('filters.silenced', function () {
      $cookieStore.put('hideSilenced', $scope.filters.silenced);
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
controllerModule.controller('StashModalCtrl', ['$scope', '$modalInstance', 'item', 'stashesService',
  function ($scope, $modalInstance, item, stashesService) {
    $scope.item = item;
    $scope.stash = {};
    $scope.stash.acknowledged = item.acknowledged;
    $scope.stash.dc = item.dc;
    $scope.stash.reason = '';
    $scope.stash.expirations = {
      '900': 900,
      '3600': 3600,
      '86400': 86400,
      'none': -1
    };
    $scope.stash.expiration = 900;
    $scope.stash.path = stashesService.construct(item);

    $scope.ok = function () {
      stashesService.submit(item, $scope.stash);
      $modalInstance.close();
    };
    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }
]);
