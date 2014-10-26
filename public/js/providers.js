'use strict';

var providerModule = angular.module('uchiwa.providers', []);

/**
* Notifications
*/
providerModule.provider('notification', function () {
  this.$get = function (toastr, toastrConfig, $cookieStore) {
    var toastrSettings = $cookieStore.get('toastrSettings');
    if(!toastrSettings) {
      toastrSettings = { 'positionClass': 'toast-bottom-right', timeOut: 5000 };
      $cookieStore.put('toastrSettings', toastrSettings);
    }
    angular.extend(toastrConfig, toastrSettings);
    return function (type, message) {
      if (type !== 'error' && type !== 'warning' && type !== 'success') {
        type = 'info';
      }
      var title = '';
      if (type === 'success') {
        var titles = ['Great!', 'All right!', 'Fantastic!', 'Excellent!', 'Good news!'];
        var rand = Math.floor((Math.random() * titles.length) + 1);
        title = titles[rand];
      }
      else if (type === 'error') {
        title = 'Oops! Something went wrong.';
      }
      toastr[type](message, title);
    };
  };
});
