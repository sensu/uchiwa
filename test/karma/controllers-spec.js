'use strict';

describe('Controller', function () {
  var $rootScope;
  var $scope;
  var createController;
  var mockNotification;
  var mockStashesService;
  var mockRoutingService;
  var mockSensuData;
  var mockVersion;

  beforeEach(module('uchiwa'));

  beforeEach(function () {
    mockNotification = jasmine.createSpy('mockNotification');
    mockStashesService = jasmine.createSpyObj('mockStashesService', ['stash', 'deleteStash']);
    mockRoutingService = jasmine.createSpyObj('mockRoutingService', ['search', 'go', 'initFilters', 'permalink', 'updateFilters']);

    mockSensuData = {
      Dc: 'abcd',
      Clients: 'efgh',
      Subscriptions: 'hijk',
      Events: 'lmno'
    };

    mockVersion = {
      uchiwa: 'x.y.z'
    };
    module(function ($provide) {
      $provide.value('notification', mockNotification);
      $provide.value('stashesService', mockStashesService);
      $provide.value('routingService', mockRoutingService);
    });
  });

  beforeEach(inject(function ($controller, _$rootScope_) {
    $rootScope = _$rootScope_;
    $scope = $rootScope.$new();
    createController = function (controllerName, properties) {
      return $controller(controllerName, _.extend({
        '$scope': $scope
      }, properties));
    };
  }));

  describe('init', function () {
    var controllerName = 'init';

    it('should call getSensu on route change success', function () {
      createController(controllerName);
      spyOn($rootScope, 'getSensu');
      $rootScope.$broadcast('$routeChangeSuccess', {});
      expect($rootScope.getSensu).toHaveBeenCalled();
    });

    //it('should create notification on messenger event', function () {
    //  var expectedType = 'success';
    //  var expectedMessage = '<strong>Success!</strong> The stash has been created.';
    //  var payload = {content: angular.toJson({type: 'success', content: expectedMessage})};
    //  createController(controllerName);

    //  socket.receive('messenger', payload);

    //  expect(mockNotification).toHaveBeenCalledWith(expectedType, expectedMessage);
    //});
  });

  describe('checks', function () {
    var controllerName = 'checks';

    beforeEach(function () {
      spyOn($scope, '$on').and.callThrough();
      createController(controllerName);
    });

    it('should have a subscribersSummary method', function () {
      expect($scope.subscribersSummary).toBeDefined();
    });

    it('should listen for the $locationChangeSuccess event', function () {
      expect($scope.$on).toHaveBeenCalledWith('$locationChangeSuccess', jasmine.any(Function));
    });
    it('should handle the $locationChangeSuccess event', function () {
      expect($scope.filters).toBeDefined();
      expect(mockRoutingService.initFilters).toHaveBeenCalled();
      $scope.$emit('$locationChangeSuccess', {});
      expect(mockRoutingService.updateFilters).toHaveBeenCalled();
    });

    describe('subscribersSummary()', function () {

      it('should join strings', function () {
        var mockArray = ['test', 'a', 'b', 'c'];
        var mockString = 'test a b c';
        expect($scope.subscribersSummary(mockArray)).toBe(mockString);
      });

    });

  });

  describe('client', function () {
    var controllerName = 'client';

    it('should have a deleteClient method', function () {
      createController(controllerName);
      expect($scope.deleteClient).toBeDefined();
    });
    it('should have a resolveEvent method', function () {
      createController(controllerName);
      expect($scope.resolveEvent).toBeDefined();
    });
    it('should have a permalink method', function () {
      createController(controllerName);
      expect($scope.permalink).toBeDefined();
    });
    it('should have a stash method', function () {
      createController(controllerName);
      expect($scope.stash).toBeDefined();
    });
  });

  describe('clients', function () {
    var controllerName = 'clients';

    it('should have a go method', function () {
      createController(controllerName);
      expect($scope.go).toBeDefined();
    });
    it('should have a stash method', function () {
      createController(controllerName);
      expect($scope.stash).toBeDefined();
    });
    it('should have a permalink method', function () {
      createController(controllerName);
      expect($scope.permalink).toBeDefined();
    });

    describe('permalink()', function () {

      it('should call routing service permalink method', function () {
        createController(controllerName);
        $scope.permalink();
        expect(mockRoutingService.permalink).toHaveBeenCalled();
      });

    })
  });

  describe('events', function () {
    var controllerName = 'events';

    describe('methods', function () {

      beforeEach(function () {
        createController(controllerName);
      });
      it('should have a go method', function () {
        expect($scope.go).toBeDefined();
      });
      it('should have a stash method', function () {
        expect($scope.stash).toBeDefined();
      });

    });
  });

  describe('info', function () {
    var controllerName = 'info';
  });

  describe('navbar', function () {
    var controllerName = 'navbar';

    it('should count events and client status on sensu', function () {
      createController(controllerName);
      var expectedClients = [
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
        },
        {
          status: 3
        }
      ];
      var expectedCriticalClients = 2;
      var expectedWarningClients = 2;
      var expectedUnknownClients = 1;
      var expectedTotalClients = 5;
      var expectedClientsStyle = 'critical';

      var expectedEvents = [
        {
          check: {
            status: 2
          }
        },
        {
          check: {
            status: 2
          }
        },
        {
          check: {
            status: 1
          }
        },
        {
          check: {
            status: 1
          }
        },
        {
          check: {
            status: 3
          }
        }
      ];
      var expectedCriticalEvents = 2;
      var expectedWarningEvents = 2;
      var expectedUnknownEvents = 1;
      var expectedTotalEvents = 5;
      var expectedEventsStyle = 'critical';

      //var payload = {Events: expectedEvents, Clients: expectedClients};
      $rootScope.events = expectedEvents;
      $rootScope.clients = expectedClients;
      $rootScope.$broadcast('sensu');

      expect($scope.clients.critical).toEqual(expectedCriticalClients);
      expect($scope.clients.warning).toEqual(expectedWarningClients);
      expect($scope.clients.unknown).toEqual(expectedUnknownClients);
      expect($scope.clients.total).toEqual(expectedTotalClients);
      expect($scope.events.critical).toEqual(expectedCriticalEvents);
      expect($scope.events.warning).toEqual(expectedWarningEvents);
      expect($scope.events.unknown).toEqual(expectedUnknownEvents);
      expect($scope.events.total).toEqual(expectedTotalEvents);
      expect($scope.clients.style).toEqual(expectedClientsStyle);
      expect($scope.events.style).toEqual(expectedEventsStyle);
    });

    it('should count unknown events and clients status on sensu', function () {
      createController(controllerName);
      var expectedClients = [
        {
          status: 3
        },
        {
          status: 3
        }
      ];
      var expectedCriticalClients = 0;
      var expectedWarningClients = 0;
      var expectedUnknownClients = 2;
      var expectedTotalClients = 2;
      var expectedClientsStyle = 'unknown';

      var expectedEvents = [
        {
          check: {
            status: 3
          }
        },
        {
          check: {
            status: 3
          }
        }
      ];
      var expectedCriticalEvents = 0;
      var expectedWarningEvents = 0;
      var expectedUnknownEvents = 2;
      var expectedTotalEvents = 2;
      var expectedEventsStyle = 'unknown';

      $rootScope.events = expectedEvents;
      $rootScope.clients = expectedClients;
      $rootScope.$broadcast('sensu');

      expect($scope.clients.critical).toEqual(expectedCriticalClients);
      expect($scope.clients.warning).toEqual(expectedWarningClients);
      expect($scope.clients.unknown).toEqual(expectedUnknownClients);
      expect($scope.clients.total).toEqual(expectedTotalClients);
      expect($scope.events.critical).toEqual(expectedCriticalEvents);
      expect($scope.events.warning).toEqual(expectedWarningEvents);
      expect($scope.events.unknown).toEqual(expectedUnknownEvents);
      expect($scope.events.total).toEqual(expectedTotalEvents);
      expect($scope.clients.style).toEqual(expectedClientsStyle);
      expect($scope.events.style).toEqual(expectedEventsStyle);
    });
  });

  describe('sidebar', function () {
    var controllerName = 'sidebar';

    it('should have a getClass method', function () {
      createController(controllerName);
      expect($scope.getClass).toBeDefined();
    });

    describe('getClass()', function () {

      it('should return selected if path matches location', function () {
        createController(controllerName, {
          '$location': {
            path: function () {
              return 'events#some-anchor';
            }
          }
        });
        expect($scope.getClass('events')).toBe('selected');
        expect($scope.getClass('clients')).toBe('');
      });

    });

  });

  describe('stashes', function () {
    var controllerName = 'stashes';

    it('should listen for the $locationChangeSuccess event', function () {
      spyOn($scope, '$on').and.callThrough();
      createController(controllerName);
      expect($scope.$on).toHaveBeenCalledWith('$locationChangeSuccess', jasmine.any(Function));
    });
    it('should handle the $locationChangeSuccess event', function () {
      createController(controllerName);
      $scope.$emit('$locationChangeSuccess', {});
      expect(mockRoutingService.updateFilters).toHaveBeenCalled();
    });

  });

  describe('settings', function () {
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
