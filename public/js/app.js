angular.module('uchiwa', [
  'uchiwa.controllers',
  'uchiwa.services',
  'uchiwa.directives',
  // Angular dependencies
  'ngRoute',
  // 3rd party dependencies
  'btford.socket-io'
]);

angular.module('uchiwa').config(['$routeProvider',
  function ($routeProvider) {
    $routeProvider
      .when('/', {templateUrl: 'partials/dashboard/index.html', controller: 'dashboard'})
      .when('/clients', {templateUrl: 'partials/client/index.html', controller: 'clients'})
      .when('/checks', {templateUrl: 'partials/check/index.html', controller: 'checks'})
      .when('/stashes', {templateUrl: 'partials/stash/index.html', controller: 'stashes'})
      .otherwise('/');
  }]);
