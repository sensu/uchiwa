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