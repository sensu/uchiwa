'use strict';

var filterModule = angular.module('uchiwa.filters', []);

filterModule.filter('encodeURIComponent', function() {
  return window.encodeURIComponent;
});

filterModule.filter('displayObject', function() {
  return function(input) {
    if(angular.isObject(input)) {
      if(input.constructor.toString().indexOf('Array') === -1) { return input; }
      var string = input.join(', ');
      return string;
    }
    else {
      return input;
    }
  };
});

filterModule.filter('filterSubscriptions', function() {
  return function(object, query) {
    if(query === '' || !object) {
      return object;
    }
    else {
      return object.filter(function (item) {
        return item.subscriptions.indexOf(query) > -1;
      });
    }
  };
});