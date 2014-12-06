'use strict';

var factoryModule = angular.module('uchiwa.factories', []);

/**
* Page title
*/
factoryModule.factory('Page', function() {
  var title = 'Uchiwa';
  return {
    title: function() { return title + ' | Uchiwa'; },
    setTitle: function(newTitle) { title = newTitle; }
  };
});

/**
* Underscore.js
*/
factoryModule.factory('underscore', function () {
  if (angular.isUndefined(window._)) {
    console.log('underscore.js is required');
  } else {
    return window._;
  }
});
