'use strict';

describe('controllers', function () {
  var $rootScope;
  var $scope;
  var createController;
  var mockNotification;
  var socket;
  var mockStashesService;

  beforeEach(module('uchiwa'));

  beforeEach(function () {
    mockNotification = jasmine.createSpy('mockNotification');
    mockStashesService = jasmine.createSpyObj('mockStashesService', ['stash']);
    module(function ($provide) {
      $provide.value('notification', mockNotification);
      $provide.value('stashesService', mockStashesService);
    });
  });

  beforeEach(inject(function ($controller, _$rootScope_, _socket_) {
    $rootScope = _$rootScope_;
    $scope = $rootScope.$new();
    socket = _socket_;
    createController = function (controllerName) {
      return $controller(controllerName, {
        '$scope': $scope
      });
    };
  }));

  describe('init controller', function () {
    var controllerName = 'init';

    it('should emit get_sensu on route change success', function () {
      createController(controllerName);
      spyOn(socket, 'emit');
      var expectedEvent = 'get_sensu';
      var expectedPayload = {};

      $rootScope.$broadcast('$routeChangeSuccess', {});

      expect(socket.emit).toHaveBeenCalledWith(expectedEvent, expectedPayload);
    });

    it('should create notification on messenger event', function () {
      var expectedType = 'success';
      var expectedMessage = '<strong>Success!</strong> The stash has been created.';
      var payload = {content: angular.toJson({type: 'success', content: expectedMessage})};
      createController(controllerName);

      socket.receive('messenger', payload);

      expect(mockNotification).toHaveBeenCalledWith(expectedType, expectedMessage);
    });
  });

  describe('checks controller', function () {
    var controllerName = 'checks';

    it('should set dc and aggregate rows on socket:sensu', function () {
      createController(controllerName);
      var expectedDataCenters = [
        {'name': 'API 1'},
        {'name': 'API 2'}
      ];
      var expectedAggregation = [
        [
          [
            {name: 'apache-running', dc: 'API 1'},
            {name: 'rabbitmq-server-running', dc: 'API 1'}
          ]
        ]
      ];
      var checks = [
        [
          {'name': 'apache-running', 'dc': 'API 1'},
          {'name': 'rabbitmq-server-running', 'dc': 'API 1'}
        ]
      ];
      var payload = {
        content: angular.toJson({dc: expectedDataCenters, checks: checks})
      };

      $rootScope.$broadcast('socket:sensu', payload);

      expect($scope.dc).toEqual(expectedDataCenters);
      expect($scope.aggregation).toEqual(expectedAggregation);
    });
  });

  describe('clients controller', function () {
    var controllerName = 'clients';

    it('should have a stash method', function () {
      createController(controllerName);

      expect($scope.stash).toBeDefined();
    });

    it('should emit get_client', function () {
      createController(controllerName);
      spyOn(socket, 'emit');
      var dataCenterName = 'foo';
      var clientName = 'bar';
      var expectedEventName = 'get_client';
      var expectedPayload = {dc: dataCenterName, client: clientName};

      $scope.getClient(dataCenterName, clientName);

      expect(socket.emit).toHaveBeenCalledWith(expectedEventName, expectedPayload);
    });

    it('should set dc and aggregate rows on socket:sensu', function () {
      createController(controllerName);
      var expectedDataCenters = [
        {'name': 'API 1'},
        {'name': 'API 2'}
      ];
      var expectedAggregation = [
        [
          [
            {name: 'sensu-client', dc: 'API 1'}
          ]
        ],
        [
          [
            {name: 'sensu-client', dc: 'API 2'}
          ]
        ]
      ];
      var clients = [
        [
          {
            'name': 'sensu-client',
            'dc': 'API 1'
          }
        ],
        [
          {
            'name': 'sensu-client',
            'dc': 'API 2'
          }
        ]
      ];
      var payload = {
        content: angular.toJson({dc: expectedDataCenters, clients: clients})
      };

      $rootScope.$broadcast('socket:sensu', payload);

      expect($scope.dc).toEqual(expectedDataCenters);
      expect($scope.aggregation).toEqual(expectedAggregation);
    });
  });

  describe('dashboard controller', function () {
    var controllerName = 'dashboard';

    it('should emit get_stats', function () {
      spyOn(socket, 'emit');
      var expectedEventName = 'get_stats';
      var expectedPayload = {};

      createController(controllerName);

      expect(socket.emit).toHaveBeenCalledWith(expectedEventName, expectedPayload);
    });

    it('should count event and client status on socket:sensu', function () {
      createController(controllerName);
      var clients = [
        [
          {
            status: 2
          },
          {
            status: 2
          },
          {
            status: 1
          },
          {
            status: 1
          }
        ]
      ];
      var expectedCriticalClients = 2;
      var expectedWarningClients = 2;

      var expectedEvents = [
        [
          {
            'check': {
              'status': 2
            }
          },
          {
            'check': {
              'status': 2
            }
          },
          {
            'check': {
              'status': 1
            }
          },
          {
            'check': {
              'status': 1
            }
          }
        ]
      ];
      var expectedCriticalEvents = 2;
      var expectedWarningEvents = 2;
      var payload = {
        content: angular.toJson({events: expectedEvents, clients: clients})
      };

      $rootScope.$broadcast('socket:sensu', payload);

      expect($scope.clients.critical).toEqual(expectedCriticalClients);
      expect($scope.clients.warning).toEqual(expectedWarningClients);
      expect($scope.events.critical).toEqual(expectedCriticalEvents);
      expect($scope.events.warning).toEqual(expectedWarningEvents);
    });

    it('should emit get_client', function () {
      createController(controllerName);
      spyOn(socket, 'emit');
      var dataCenterName = 'foo';
      var clientName = 'bar';
      var expectedEventName = 'get_client';
      var expectedPayload = {dc: dataCenterName, client: clientName};

      $scope.getClient(dataCenterName, clientName);

      expect(socket.emit).toHaveBeenCalledWith(expectedEventName, expectedPayload);
    });
  });

  describe('stashes controller', function () {
    var controllerName = 'stashes';

    it('should set dc and aggregate rows on socket:sensu', function () {
      createController(controllerName);
      var expectedDataCenters = [
        {'name': 'API 1'},
        {'name': 'API 2'}
      ];
      var expectedAggregation = [
        [
          [
            {name: 'stash 1', dc: 'API 1'}
          ]
        ],
        [
          [
            {name: 'stash 2', dc: 'API 1'}
          ]
        ]
      ];
      var stashes = [
        [
          {name: 'stash 1', dc: 'API 1'}
        ],
        [
          {name: 'stash 2', dc: 'API 1'}
        ]
      ];
      var payload = {content: angular.toJson({dc: expectedDataCenters, stashes: stashes})};

      $rootScope.$broadcast('socket:sensu', payload);

      expect($scope.dc).toEqual(expectedDataCenters);
      expect($scope.aggregation).toEqual(expectedAggregation);
    });

    it("should delete stash", function () {
      createController(controllerName);
      var expectedDcName = 'API 1';
      var stashToDelete = {name: 'stash'};
      var stashesForDc = [[[stashToDelete]]];
      $scope.sensu = { dc: [{name: expectedDcName}], stashes: stashesForDc };

      $scope.deleteStash('API 1', stashToDelete, 0);

      expect($scope.sensu.stashes[0][0]).not.toContain(stashToDelete);
      expect(mockStashesService.stash).toHaveBeenCalledWith(expectedDcName, stashToDelete);
    });
  });

  describe('settings controller', function () {
    var controllerName = 'settings';

    it("should emit a theme:changed event when the current theme changes", function () {
      createController(controllerName);
      var expectedTheme = 'foo theme';
      var expectedEvent = 'theme:changed';
      spyOn($scope, '$emit');

      $scope.currentTheme = expectedTheme;

      $scope.$apply();
      expect($scope.$emit).toHaveBeenCalledWith(expectedEvent, expectedTheme);
    });
  });
});