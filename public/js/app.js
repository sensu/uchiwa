angular.module('uchiwa', [
  'uchiwa.controllers',
  'uchiwa.constants',
  'uchiwa.services',
  'uchiwa.directives',
  // Angular dependencies
  'ngCookies',
  'ngRoute',
  // 3rd party dependencies
  'btford.socket-io'
]);

angular.module('uchiwa').config(['$routeProvider', 'notificationProvider',
  function ($routeProvider, notificationProvider) {
    $routeProvider
      .when('/', {templateUrl: 'partials/dashboard/index.html', controller: 'dashboard'})
      .when('/clients', {templateUrl: 'partials/client/index.html', controller: 'clients'})
      .when('/checks', {templateUrl: 'partials/check/index.html', controller: 'checks'})
      .when('/stashes', {templateUrl: 'partials/stash/index.html', controller: 'stashes'})
      .when('/settings', {templateUrl: 'partials/settings/edit.html', controller: 'settings'})
      .otherwise('/');
    notificationProvider.setOptions({
      "positionClass": "toast-bottom-right"
    });
  }]);
