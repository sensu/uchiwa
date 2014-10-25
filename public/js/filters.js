'use strict';

var filterModule = angular.module('uchiwa.filters', []);

filterModule.filter('arrayToString', function() {
  return function(array) {
    if (!array) { return ''; }
    if (array.constructor.toString().indexOf('Array') === -1) { return array; }
    return array.join(' ');
  };
});

filterModule.filter('buildStashes', function() {
  return function(stashes) {
    if (Object.prototype.toString.call(stashes) !== '[object Array]') {
      return stashes;
    }
    angular.forEach(stashes, function(stash) {
      var path = stash.path.split('/');
      stash.client = path[1] || null;
      stash.check = path[2] || null;
    });
    return stashes;
  };
});

filterModule.filter('buildEvents', function() {
  return function(events) {
    if (Object.prototype.toString.call(events) !== '[object Array]') {
      return events;
    }
    angular.forEach(events, function(event) {
      event.sourceName = event.check.source || event.client.name;
    });
    return events;
  };
});

filterModule.filter('displayObject', function() {
  return function(input) {
    if(angular.isObject(input)) {
      if(input.constructor.toString().indexOf('Array') === -1) { return input; }
      return input.join(', ');
    }
    else {
      return input;
    }
  };
});

filterModule.filter('encodeURIComponent', function() {
  return window.encodeURIComponent;
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

filterModule.filter('getAckClass', function() {
  return function(isAcknowledged) {
    return (isAcknowledged) ? 'fa-volume-off' : 'fa-volume-up';
  };
});

filterModule.filter('getExpireTimestamp', ['$filter', 'settings', function ($filter, settings) {
  return function(expire, timestamp) {
    if (isNaN(timestamp) || isNaN(expire)) {
      return 'Unknown';
    }
    if (expire === -1) {
      return 'Never';
    }
    var expiration = (expire + timestamp) * 1000;
    return $filter('date')(expiration, settings.date);
  };
}]);

filterModule.filter('getStatusClass', function() {
  return function(status) {
    switch(status) {
      case 0:
        return 'success';
      case 1:
        return 'warning';
      case 2:
        return 'critical';
      default:
        return 'unknown';
    }
  };
});

filterModule.filter('getTimestamp', ['$filter', 'settings', function ($filter, settings) {
  return function(timestamp) {
    if (isNaN(timestamp) || timestamp.toString().length !== 10) {
      return timestamp;
    }
    timestamp = timestamp * 1000;
    return $filter('date')(timestamp, settings.date);
  };
}]);

filterModule.filter('setMissingProperty', function() {
  return function(property) {
    return property || false;
  };
});

filterModule.filter('richOutput', ['$filter', function($filter) {
  return function(text) {
    var output = '';
    if(typeof text === 'object') {
      if (text instanceof Array) {
        output = text.join(', ');
      } else {
        var code = hljs.highlight('json', angular.toJson(text, true)).value;
        output = '<pre class=\"hljs\">' + code + '</pre>';
      }
    } else if (typeof text === 'number') {
      output = text.toString();
    } else {
      var linkified = $filter('linky')(text, '_blank');
      output = $filter('imagey')(linkified);
    }
    return output;
  };
}]);
