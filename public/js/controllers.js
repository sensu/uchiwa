'use strict';

var controllerModule = angular.module('uchiwa.controllers', []);

/**
 * Init
 */
controllerModule.controller('init', ['$scope', 'notification', 'socket', 'Page',
  function ($scope, notification, socket, Page) {
    $scope.Page = Page;
    $scope.$on('$routeChangeSuccess', function () {
      socket.emit('get_sensu', {});
    });

    socket.on('messenger', function (data) {
      if (angular.isDefined(data.content)) {
        var message = angular.fromJson(data.content);
        notification(message.type, message.content);
      }
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

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.dc = sensu.dc;
      $scope.checks = sensu.checks;
    });

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
controllerModule.controller('client', ['$scope', '$routeParams', 'socket', 'clientsService', 'routingService', 'Page',
  function ($scope, $routeParams, socket, clientsService, routingService, Page) {

    $scope.predicate = '-last_status';

    // Retrieve client
    $scope.clientId = decodeURI($routeParams.clientId);
    $scope.dcId = decodeURI($routeParams.dcId);
    $scope.pull = function() {
      socket.emit('get_client', {dc: $scope.dcId, client: $scope.clientId});
    };

    $scope.pull();
    var timer = setInterval($scope.pull, 10000);

    // Socket.IO
    $scope.$on('socket:client', function (event, data) {
      if(!$scope.dropdown.isopen) {
        $scope.client = angular.fromJson(data.content);
      }
      $scope.pageHeaderText = $scope.client.name;

      // Retrieve check
      $scope.requestedCheck = decodeURI($routeParams.check);
      $scope.selectedCheck = findCheck($scope.requestedCheck);

      // Set page title
      if(angular.isDefined($scope.selectedCheck)) {
        Page.setTitle($scope.requestedCheck + ' - ' + $scope.client.name);
      }
      else {
        Page.setTitle($scope.client.name);
      }
    });

    // Helpers
    $scope.toggled = function(e) {
      var event = e || window.event;
      event.stopPropagation();

      $scope.dropdown.isopen = !$scope.dropdown.isopen;
    };
    $scope.dropdown = {
      isopen: false
    };
    $scope.silenceOptions = [
      {key: '15 minutes', value: 900},
      {key: '1 hour', value: 3600},
      {key: '24 hours', value: 86400},
      {key: 'Never', value: -1},
    ];
    var findCheck = function(id){
      return $scope.client.history.filter(function (item) {
        return item.check === id;
      })[0];
    };

    // Listeners
    $scope.$on('$routeUpdate', function(){
      // Update check
      $scope.requestedCheck = decodeURI($routeParams.check);
      $scope.selectedCheck = findCheck($scope.requestedCheck);
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

    // Sanitize - only display useful information
    /* jshint ignore:start */
    var clientWhitelist = [ 'dc', 'events', 'eventsSummary', 'history', 'isSilenced', 'lastCheck', 'silenceIcon', 'status', 'timestamp', 'style' ];
    var checkWhitelist = [ 'dc', 'hasSubscribers', 'name'];
    var eventWhitelist = [ 'command', 'executed', 'handlers', 'hasSubscribers', 'history', 'interval', 'issued', 'name', 'status', 'standalone', 'subscribers' ];
    $scope.sanitizeObject = function(type, key){
      return eval(type + 'Whitelist').indexOf(key) === -1;
    };
    /* jshint ignore:end */

    // Services
    $scope.remove = clientsService.remove;
    $scope.resolve = clientsService.resolve;
    $scope.permalink = routingService.permalink;
    $scope.stash = clientsService.stash;
  }
]);

/**
 * Clients
 */
controllerModule.controller('clients', ['$scope', '$routeParams', 'socket', 'clientsService', 'routingService', 'Page',
  function ($scope, $routeParams, socket, clientsService, routingService, Page) {
    Page.setTitle('Clients');
    $scope.pageHeaderText = 'Clients';
    $scope.predicate = '-status';

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.dc = sensu.dc;
      $scope.subscriptions = sensu.subscriptions;
      if(!$scope.dropdown.isopen) {
        $scope.clients = sensu.clients;
      }
    });

    // Helpers
    $scope.getClient = function (dcName, clientName) {
      socket.emit('get_client', {dc: dcName, client: clientName});
    };

    // Routing
    $scope.filters = {};
    routingService.initFilters($routeParams, $scope.filters, ['dc', 'subscription', 'limit', 'q']);
    $scope.$on('$locationChangeSuccess', function(){
      routingService.updateFilters($routeParams, $scope.filters);
    });

    // Services
    $scope.go = routingService.go;
    $scope.stash = clientsService.stash;
    $scope.permalink = routingService.permalink;

    $scope.toggled = function(e) {
      var event = e || window.event;
      event.stopPropagation();
      $scope.dropdown.isopen = !$scope.dropdown.isopen;
    };
    $scope.dropdown = {
      isopen: false
    };
    $scope.silenceOptions = [
      {key: '15 minutes', value: 900},
      {key: '1 hour', value: 3600},
      {key: '24 hours', value: 86400},
      {key: 'Never', value: -1},
    ];
  }
]);

/**
 * Events
 */
controllerModule.controller('events', ['$scope', '$routeParams', 'socket', 'eventsService', 'routingService', 'Page',
  function ($scope, $routeParams, socket, eventsService, routingService, Page) {
    Page.setTitle('Events');
    $scope.pageHeaderText = 'Events';
    $scope.predicate = '-check.status';

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.clients = sensu.clients;
      $scope.dc = sensu.dc;
      $scope.subscriptions = sensu.subscriptions;
      if(!$scope.dropdown.isopen) {
        $scope.events = sensu.events;
      }
    });

    // Helpers
    $scope.getClient = function (dcName, clientName) {
      socket.emit('get_client', {dc: dcName, client: clientName});
    };

    // Routing
    $scope.filters = {};
    routingService.initFilters($routeParams, $scope.filters, ['dc', 'limit', 'q']);
    $scope.$on('$locationChangeSuccess', function(){
      routingService.updateFilters($routeParams, $scope.filters);
    });

    // Services
    $scope.go = routingService.go;
    $scope.stash = eventsService.stash;
    $scope.permalink = routingService.permalink;

    $scope.toggled = function(e) {
      var event = e || window.event;
      event.stopPropagation();
      $scope.dropdown.isopen = !$scope.dropdown.isopen;
    };

    $scope.dropdown = {
      isopen: false
    };
    $scope.silenceOptions = [
      {key: '15 minutes', value: 900},
      {key: '1 hour', value: 3600},
      {key: '24 hours', value: 86400},
      {key: 'Never', value: -1},
    ];

  }
]);

/**
 * Info
 */
controllerModule.controller('info', ['$scope', 'socket', 'version', 'Page',
  function ($scope, socket, version, Page) {
    $scope.pageHeaderText = 'Info';
    $scope.uchiwa = {};
    Page.setTitle('Info');

    // Socket.IO
    socket.emit('get_info', {});

    $scope.$on('socket:info', function (event, data) {
      var config = angular.fromJson(data.content);
      $scope.uchiwa.config = JSON.stringify(config, null, 2);
      $scope.uchiwa.version = version.uchiwa;
    });

    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.dc = sensu.dc;
    });

  }
]);

/**
 * Navbar
 */
controllerModule.controller('navbar', ['$scope', 'routingService',
  function ($scope, routingService) {

    // Services
    $scope.go = routingService.go;

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {

      var sensu = angular.fromJson(data.content);
      $scope.checks = sensu.checks;
      $scope.clients = sensu.clients;
      $scope.dc = sensu.dc;
      $scope.events = sensu.events;
      $scope.stashes = sensu.stashes;

      // Badges count
      $scope.countStatuses = function (collection, getStatusCode) {
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

      $scope.countStatuses($scope.clients, function (item) {
        return item.status;
      });
      $scope.countStatuses($scope.events, function (item) {
        return item.check.status;
      });
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
controllerModule.controller('stashes', ['$scope', '$routeParams', 'socket', 'routingService', 'stashesService', 'Page',
  function ($scope, $routeParams, socket, routingService, stashesService, Page) {
    Page.setTitle('Stashes');
    $scope.pageHeaderText = 'Stashes';
    $scope.sensu = {};
    $scope.predicate = 'client';

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      $scope.sensu = angular.fromJson(data.content);
      $scope.dc = $scope.sensu.dc;
      $scope.stashes = $scope.sensu.stashes;
    });

    // Helpers
    $scope.deleteStash = function (dcName, stash, index) {
      stashesService.stash(dcName, stash);
      $scope.stashes.splice(index, 1);
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