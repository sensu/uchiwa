var controllerModule = angular.module('uchiwa.controllers', []);

/**
 * Init
 */
controllerModule.controller('init', ['$scope', 'notification', 'socket',
  function ($scope, notification, socket) {
    $scope.$on('$routeChangeSuccess', function () {
      socket.emit('get_sensu', {});
    });

    socket.on('messenger', function (data) {
      if (angular.isDefined(data.content)) {
        var message = angular.fromJson(data.content);
        notification(message.type, message.content);
      }
    });
    $scope.filters = {text: ''};
  }
]);

/**
 * Checks
 */
controllerModule.controller('checks', ['$scope', 'utilityService', 'toggleService',
  function ($scope, utilityService, toggleService) {
    $scope.pageHeaderText = 'Checks';
    $scope.dcItem = 'checks';

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = utilityService.getRows(sensu.checks, 2);
    });

    // Toggle system
    $scope.toggle = toggleService.toggle;
    $scope.toggleOn = toggleService.toggleOn;
    $scope.showOnly = toggleService.showOnly;
    $scope.showAll = toggleService.showAll;
  }
]);

/**
 * Client
 */
controllerModule.controller('client', ['$scope', '$location', 'socket', 'clientsService', 'toggleService', 'toggleClientService',
  function ($scope, $location, socket, clientsService, toggleService, toggleClientService) {
    var timer = setInterval(function () {
      if ($('#client-details').data('bs.modal')) {
        socket.emit('get_client', {dc: $scope.client.dc, client: $scope.client.name});
      }
    }, 10000);
    $scope.$on('socket:client', function (event, data) {
      $scope.client = angular.fromJson(data.content);
    });
    $scope.stash = clientsService.stash;
    $scope.resolve = clientsService.resolve;

    $scope.delete = function (dcName, clientName) {
      clientsService.delete(dcName, clientName);
      $('#client-details').modal('hide');
    };
    $('#client-details').on('hide.bs.modal', function () {
      $scope.client = {name: 'Loading...'};
      $scope.toggle = toggleService.toggle;
      clearInterval(timer);
    });

    // Keep track of collapsed check details
    $scope.toggleClient = toggleClientService.toggle;
    $scope.toggleClientActive = toggleClientService.toggleOn;
  }
]);

/**
 * Clients
 */
controllerModule.controller('clients', ['$scope', 'socket', 'clientsService', 'utilityService', 'toggleService',
  function ($scope, socket, clientsService, utilityService, toggleService) {
    $scope.pageHeaderText = 'Clients';
    $scope.dcItem = 'clients';

    // Helpers
    $scope.stash = clientsService.stash;
    $scope.getClient = function (dcName, clientName) {
      socket.emit('get_client', {dc: dcName, client: clientName});
    };

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = utilityService.getRows(sensu.clients, 3);
    });

    // Toggle system
    $scope.toggle = toggleService.toggle;
    $scope.toggleOn = toggleService.toggleOn;
    $scope.showOnly = toggleService.showOnly;
    $scope.showAll = toggleService.showAll;
  }
]);

/**
 * Dashboard
 */
controllerModule.controller('dashboard', ['$rootScope', '$scope', 'socket', 'eventsService', 'utilityService', 'toggleService',
  function ($rootScope, $scope, socket, eventsService, utilityService, toggleService) {
    $scope.pageHeaderText = 'Events';
    $scope.dcItem = 'events';
    $scope.statChartConfig = {
      data: [],
      xkey: 'y',
      ykeys: ['e', 's'],
      labels: ['Events', 'Stashes'],
      lineColors: ['#2CA7E5', '#F9CD65'],
      hideHover: 'auto',
      pointSize: 0,
      fillOpacity: 1,
      gridTextColor: '#fff',
      gridTextFamily: '"Lato", sans-serif',
      gridTextWeight: 700,
      grid: false,
      lineWidth: 4,
      axes: true,
      behaveLikeLine: true
    };
    socket.emit('get_stats', {});
    $scope.$on('socket:stats', function (event, data) {
      $scope.statChartConfig.data = angular.fromJson(data.content);
    });

    $scope.getStatusClass = function (statuses) {
      if (angular.isDefined(statuses)) {
        return statuses.critical > 0 ? 'critical' : statuses.warning > 0 ? 'warning' : 'success';
      }
    };

    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.clients = sensu.clients;
      $scope.events = sensu.events;
      $scope.countStatuses = function (collection, getStatusCode) {
        var criticals = 0;
        var warnings = 0;

        angular.forEach(collection, function (item) {
          criticals += item.filter(function (item) {
            return getStatusCode(item) === 2;
          }).length;
          warnings += item.filter(function (item) {
            return getStatusCode(item) === 1;
          }).length;
        });

        collection.warning = warnings;
        collection.critical = criticals;
        collection.total = criticals + warnings;
      };

      $scope.countStatuses($scope.clients, function (item) {
        return item.status;
      });
      $scope.countStatuses($scope.events, function (item) {
        return item.check.status;
      });
    });

    // Helpers
    $scope.stash = eventsService.stash;
    $scope.getClient = function (dcName, clientName) {
      socket.emit('get_client', {dc: dcName, client: clientName});
    };

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      var sensu = angular.fromJson(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = utilityService.getRows(sensu.events, 3);
    });

    // Toggle system
    $scope.toggle = toggleService.toggle;
    $scope.toggleOn = toggleService.toggleOn;
    $scope.showOnly = toggleService.showOnly;
    $scope.showAll = toggleService.showAll;
  }
]);

/**
 * Stashes
 */
controllerModule.controller('stashes', ['$scope', 'socket', 'stashesService', 'utilityService', 'toggleService',
  function ($scope, socket, stashesService, utilityService, toggleService) {
    $scope.pageHeaderText = 'Stashes';
    $scope.dcItem = 'stashes';
    $scope.sensu = {};

    // Socket.IO
    $scope.$on('socket:sensu', function (event, data) {
      $scope.sensu = angular.fromJson(data.content);
      $scope.dc = $scope.sensu.dc;
      $scope.aggregation = utilityService.getRows($scope.sensu.stashes, 3);
    });

    $scope.deleteStash = function (dcName, stash, index) {
      stashesService.stash(dcName, stash);

      // Remove stash from $scope
      var dcPosition = $scope.sensu.dc.map(function (dc) {
        return dc.name;
      }).indexOf(dcName);
      var dcStashes = $scope.sensu.stashes[dcPosition];
      dcStashes[0].splice(index, 1);
    };

    // Toggle system
    $scope.toggle = toggleService.toggle;
    $scope.toggleOn = toggleService.toggleOn;
    $scope.showOnly = toggleService.showOnly;
    $scope.showAll = toggleService.showAll;
  }
]);

/**
 * Settings
 */
controllerModule.controller('settings', ['$cookies', '$scope',
  function ($cookies, $scope) {
    $scope.pageHeaderText = 'Settings';
    $scope.$watch('currentTheme', function (theme) {
      $scope.$emit('theme:changed', theme);
    });
  }
]);
