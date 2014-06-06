var controllerModule = angular.module('uchiwa.controllers', []);

/**
 * Init
 */
controllerModule.controller('init', ['$scope', 'socket',
  function($scope, socket) {
    socket.emit('get_sensu', {});
  }
]);

/**
 * Checks
 */
controllerModule.controller('checks', ['$scope', 'socket',
  function($scope, socket) {

    // Helpers
    $scope.splitArray = function(array, n) { // Divide an array into 'n' arrays 
      var arrays = [];
      var i,j,temparray, chunk = n;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          arrays.push(temparray);
      }
      return arrays;
    };
    $scope.getRows = function(array, n) { // Get rows for each DC
      _.each(array, function(element,index,list){
        list[index] = $scope.splitArray(element,n);
      });
      return array;
    };

    // Socket.IO
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = $scope.getRows(sensu.checks, 2);
    });

    // Toggle system
    $scope.toggle = {};
    $scope.toggleOn = function (index) {
      if(typeof $scope.toggle[index] === "undefined") $scope.toggle[index] = {hidden: false};
      $scope.toggle[index].hidden = !$scope.toggle[index].hidden;
    };
    $scope.showOnly = function (index, dc) {
      _.each(dc, function(datacenter, i){
        if(i == index) return $scope.toggle[index] = {hidden: false};
        $scope.toggle[i] = {hidden: true};
      });
    };
    $scope.showAll = function (dc) {
      _.each(dc, function(datacenter, i){
        $scope.toggle[i] = {hidden: false};
      });
    };
  }
]);

/**
 * Client
 */
controllerModule.controller('client', ['$scope', '$location', 'socket', 'clientsService',
  function($scope, $location, socket, clientsService) {
    var timer = setInterval(function(){
      if($("#client-details").data('bs.modal')){
        socket.emit('get_client', {dc: $scope.client.dc, client: $scope.client.name});
      }
    }, 10000);
    $scope.$on('socket:client', function(event, data) {
      var client = JSON.parse(data.content);
      $scope.client = client;
    });
    $scope.stash = function(e, dcName, client, check){
      clientsService.stash(e, dcName, client, check);
    };
    $scope.resolve = function(e, dcName, client, check){
      clientsService.resolve(e, dcName, client, check);
    };
    $scope.delete = function(dcName, clientName){
      clientsService.delete(dcName, clientName);
      $('#client-details').modal('hide');
    };
    $('#client-details').on('hide.bs.modal', function () {
      $scope.client = {name: "Loading..."};
      $scope.toggle = {};
      clearInterval(timer);
    });

    // Keep track of collapsed check details
    $scope.toggle = {};
    $scope.toggleActive = function (index) {
      if(typeof $scope.toggle[index] === "undefined") $scope.toggle[index] = {hidden: true};
      $scope.toggle[index].hidden = !$scope.toggle[index].hidden;
    };
  }
]);

/**
 * Clients
 */
controllerModule.controller('clients', ['$scope', 'socket', 'clientsService',
  function($scope, socket, clientsService) {

    // Helpers
    $scope.stash = function(e, dcName, client){
      clientsService.stash(e, dcName, client);
    };
    $scope.getClient = function(dcName, clientName){
      socket.emit('get_client', {dc: dcName, client: clientName});
    }
    $scope.splitArray = function(array, n) { // Divide an array into 'n' arrays 
      var arrays = [];
      var i,j,temparray, chunk = n;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          arrays.push(temparray);
      }
      return arrays;
    };
    $scope.getRows = function(array, n) { // Get rows for each DC
      _.each(array, function(element,index,list){
        list[index] = $scope.splitArray(element,n);
      });
      return array;
    };

    // Socket.IO
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = $scope.getRows(sensu.clients, 3);
    });

    // Toggle system
    $scope.toggle = {};
    $scope.toggleOn = function (index) {
      if(typeof $scope.toggle[index] === "undefined") $scope.toggle[index] = {hidden: false};
      $scope.toggle[index].hidden = !$scope.toggle[index].hidden;
    };
    $scope.showOnly = function (index, dc) {
      _.each(dc, function(datacenter, i){
        if(i == index) return $scope.toggle[index] = {hidden: false};
        $scope.toggle[i] = {hidden: true};
      });
    };
    $scope.showAll = function (dc) {
      _.each(dc, function(datacenter, i){
        $scope.toggle[i] = {hidden: false};
      });
    };
  }
]);

/**
 * Dashboard
 */
controllerModule.controller('dashboard', ['$scope', 'socket',
  function($scope, socket) {
    socket.emit('get_stats', {});
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.clients = sensu.clients;
      $scope.events = sensu.events;
      $scope.countEvents = function() {
        var criticals = 0;
        var warnings = 0;
        _.each($scope.events, function(element){
          criticals += element.filter(function (e){ return e.check.status === 2 }).length;
          warnings += element.filter(function (e){ return e.check.status === 1 }).length;
        });

        // Display counts
        $scope.events.warning = warnings;
        $scope.events.critical = criticals;
        $scope.events.total = criticals + warnings;

        // Return style
        return (criticals > 0) ? "critical" : (warnings > 0) ? "warning" : "success";
      };
      $scope.countClients = function() {
        var criticals = 0;
        var warnings = 0;
        var total = 0;
        _.each($scope.clients, function(element){
          criticals += element.filter(function (e){ return e.status === 2 }).length;
          warnings += element.filter(function (e){ return e.status === 1 }).length;
          total += element.length;
        });

        // Display counts
        $scope.clients.warning = warnings;
        $scope.clients.critical = criticals;
        $scope.clients.total = total;

        // Return style
        return (criticals > 0) ? "critical" : (warnings > 0) ? "warning" : "success";
      };
      $scope.clientsStyle = function() {
        return 0;
        var criticals = $scope.clients.filter(function (e){ return e.status === 2 }).length;
        if(criticals > 0) return "critical";
        var warnings = $scope.clients.filter(function (e){ return e.status === 1 }).length;
        return (warnings > 0) ? "warning" : "success";
      };
    });
  }
]);

/**
 * Events
 */
controllerModule.controller('events', ['$scope', 'socket', 'eventsService',
  function($scope, socket, eventsService) {

    // Helpers
    $scope.stash = function(e, dcName, event){
      eventsService.stash(e, dcName, event);
    };
    $scope.getClient = function(dcName, clientName){
      socket.emit('get_client', {dc: dcName, client: clientName});
    }
    $scope.getDcStatus = function(index, clients){
      dcService.status(index, clients);
    }
    $scope.splitArray = function(array, n) { // Divide an array into 'n' arrays 
      var arrays = [];
      var i,j,temparray, chunk = n;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          arrays.push(temparray);
      }
      return arrays;
    };
    $scope.getRows = function(array, n) { // Get rows for each DC
      _.each(array, function(element,index,list){
        list[index] = $scope.splitArray(element,n);
      });
      return array;
    };

    // Socket.IO
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = $scope.getRows(sensu.events, 3);
    });

    // Toggle system
    $scope.toggle = {};
    $scope.toggleOn = function (index) {
      if(typeof $scope.toggle[index] === "undefined") $scope.toggle[index] = {hidden: false};
      $scope.toggle[index].hidden = !$scope.toggle[index].hidden;
    };
    $scope.showOnly = function (index, dc) {
      _.each(dc, function(datacenter, i){
        if(i == index) return $scope.toggle[index] = {hidden: false};
        $scope.toggle[i] = {hidden: true};
      });
    };
    $scope.showAll = function (dc) {
      _.each(dc, function(datacenter, i){
        $scope.toggle[i] = {hidden: false};
      });
    };
  }
]);

/**
 * Stashes
 */
controllerModule.controller('stashes', ['$scope', 'socket', 'stashesService',
  function($scope, socket, stashesService) {
    
    // Helpers
    $scope.splitArray = function(array, n) { // Divide an array into 'n' arrays 
      var arrays = [];
      var i,j,temparray, chunk = n;
      for (i=0,j=array.length; i<j; i+=chunk) {
          temparray = array.slice(i, i+chunk);
          arrays.push(temparray);
      }
      return arrays;
    };
    $scope.getRows = function(array, n) { // Get rows for each DC
      _.each(array, function(element,index,list){
        list[index] = $scope.splitArray(element,n);
      });
      return array;
    };

    // Socket.IO
    $scope.$on('socket:sensu', function(event, data) {
      var sensu = JSON.parse(data.content);
      $scope.dc = sensu.dc;
      $scope.aggregation = $scope.getRows(sensu.stashes, 3);

      $scope.deleteStash = function(dcName, stash, index){
        stashesService.stash(dcName, stash);

        // Remove stash from $scope
        var dcPosition = sensu.dc.indexOf(dcName);
        var dcStashes = sensu.stashes[dcPosition];    
        dcStashes[0].splice(index, 1);
      };
    }); 

    // Toggle system
    $scope.toggle = {};
    $scope.toggleOn = function (index) {
      if(typeof $scope.toggle[index] === "undefined") $scope.toggle[index] = {hidden: false};
      $scope.toggle[index].hidden = !$scope.toggle[index].hidden;
    };
    $scope.showOnly = function (index, dc) {
      _.each(dc, function(datacenter, i){
        if(i == index) return $scope.toggle[index] = {hidden: false};
        $scope.toggle[i] = {hidden: true};
      });
    };
    $scope.showAll = function (dc) {
      _.each(dc, function(datacenter, i){
        $scope.toggle[i] = {hidden: false};
      });
    };
  }
]);