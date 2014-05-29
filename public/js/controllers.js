var controllerModule = angular.module('uchiwa.controllers', []);

/**
 * Checks
 */
controllerModule.controller('checks', ['$scope', 'socket',
  function($scope, socket) {
    $scope.getRows = function(array, columns) {
      var rows = [];
      var i,j,temparray, chunk = columns;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          rows.push(temparray);
      }
      return rows;
    };
    socket.emit('get_sensu', {});
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.rows = $scope.getRows(sensu.checks,2);
    });
  }
])

/**
 * Client
 */
controllerModule.controller('client', ['$scope', 'socket', 'clientsService',
  function($scope, socket, clientsService) {
    var timer = setInterval(function(){
      if($("#client-details").data('bs.modal')){
        socket.emit('get_client', {name: $scope.client.name});
      }
    }, 10000);
    $scope.$on('socket:client', function(event, data) {
      var client = JSON.parse(data.content);
      $scope.client = client;
    });
    $scope.stash = function(e, client, check){
      clientsService.stash(e, client, check);
    };
    $scope.resolve = function(e, client, check){
      clientsService.resolve(e, client, check);
    };
    $('#client-details').on('hide.bs.modal', function () {
      $scope.client = {name: "Loading..."};
      clearInterval(timer);
    });
  }
]);

/**
 * Clients
 */
controllerModule.controller('clients', ['$scope', 'socket', 'clientsService',
  function($scope, socket, clientsService) {
    $scope.stash = function(e, client){
      clientsService.stash(e, client);
    };
    $scope.getClient = function(clientName){
      socket.emit('get_client', {name: clientName});
    }
    $scope.getRows = function(array, columns) {
      var rows = [];
      var i,j,temparray, chunk = columns;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          rows.push(temparray);
      }
      return rows;
    };
    socket.emit('get_sensu', {});
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.rows = $scope.getRows(sensu.clients,4);
    });
  }
]);

/**
 * Dashboard
 */
controllerModule.controller('dashboard', ['$scope', 'socket',
  function($scope, socket) {
    socket.emit('get_sensu', {});
    socket.emit('get_stats', {});
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.clients = sensu.clients;
      $scope.events = sensu.events;
      $scope.eventsStyle = function() {
        var criticals = $scope.events.filter(function (e){ return e.status === 2 }).length;
        if(criticals > 0) return "critical";
        var warnings = $scope.events.filter(function (e){ return e.status === 2 }).length;
        return (warnings > 0) ? "warning" : "success";
      };
      $scope.clientsStyle = function() {
        var criticals = $scope.clients.filter(function (e){ return e.status === 2 }).length;
        if(criticals > 0) return "critical";
        var warnings = $scope.clients.filter(function (e){ return e.status === 2 }).length;
        return (warnings > 0) ? "warning" : "success";
      };
      $scope.countClients = function(status) {
        if(status == 0) return $scope.clients.length;
        return ($scope.clients.filter(function (e){ return e.status === status }).length);
      };
      $scope.countEvents = function(status) {
        if(status == 0) return $scope.events.length;
        var count = $scope.events.filter(function (e){ return e.status === status }).length;
        return count;
      };
    });

  }
]);

/**
 * Events
 */
controllerModule.controller('events', ['$scope', 'socket', 'eventsService',
  function($scope, socket, eventsService) {
    $scope.getRows = function(array, columns) {
      var rows = [];
      var i,j,temparray, chunk = columns;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          rows.push(temparray);
      }
      return rows;
    };
    socket.emit('get_sensu', {});
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.rows = $scope.getRows(sensu.events,4);
    });
    $scope.stash = function(e, event){
      eventsService.stash(e, event);
    };
    $scope.getClient = function(clientName){
      socket.emit('get_client', {name: clientName});
    }
  }
]);

/**
 * Stashes
 */
controllerModule.controller('stashes', ['$scope', 'socket', 'stashesService',
  function($scope, socket, stashesService) {
    $scope.getRows = function(array, columns) {
      var rows = [];
      var i,j,temparray, chunk = columns;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          rows.push(temparray);
      }
      return rows;
    };
    socket.emit('get_sensu', {});
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.rows = $scope.getRows(sensu.stashes,3);
      $scope.deleteStash = function(stash, index){
        stashesService.stash(stash);
        sensu.stashes.splice(index, 1);
        $scope.rows = $scope.getRows(sensu.stashes,3);
      };
    }); 
  }
]);