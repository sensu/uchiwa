'use strict';

var filterModule = angular.module('uchiwa.filters', []);

filterModule.filter('arrayLength', function() {
  return function(array) {
    if (!array) { return 0; }
    if (array.constructor.toString().indexOf('Array') === -1) { return 0; }
    return array.length;
  };
});

filterModule.filter('arrayToString', function() {
  return function(array) {
    if (!array) { return ''; }
    if (array.constructor.toString().indexOf('Array') === -1) { return array; }
    return array.join(' ');
  };
});

filterModule.filter('buildEvents', function() {
  return function(events) {
    if (Object.prototype.toString.call(events) !== '[object Array]') {
      return events;
    }
    angular.forEach(events, function(event) {
      if (typeof(event.check) === 'undefined' && typeof(event.client) === 'undefined') {
        event.sourceName = 'unknown';
        return true;
      }
      else if (typeof(event.check) === 'undefined') {
        event.check = {};
      }
      event.sourceName = event.check.source || event.client.name;
    });
    return events;
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
  return function(stash) {
    if (isNaN(stash.expire)) {
      return 'Unknown';
    }
    if (stash.expire === -1) {
      return 'Never';
    }
    var expiration = (stash.content.timestamp + stash.expire) * 1000;
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

filterModule.filter('hideSilenced', function() {
  return function(events, hideSilenced) {
    if (Object.prototype.toString.call(events) !== '[object Array]') {
      return events;
    }
    if (events && hideSilenced) {
      return events.filter(function (item) {
        return item.acknowledged === false;
      });
    }
    return events;
  };
});

filterModule.filter('hideOccurrences', function() {
  return function(events, hideOccurrences) {
    if (Object.prototype.toString.call(events) !== '[object Array]') {
      return events;
    }
    if (events && hideOccurrences) {
      return events.filter(function (item) {
        if (('occurrences' in item.check) && !isNaN(item.check.occurrences)) {
          return item.occurrences >= item.check.occurrences;
        } else {
          return true;
        }
      });
    }
    return events;
  };
});

filterModule.filter('imagey', function() {
  return function(url) {
    if (!url) {
      return url;
    }
    var IMG_URL_REGEX = /(href=['"]?)?https?:\/\/(?:[0-9a-zA-Z_\-\.]+)\/(?:[^'"]+)\.(?:jpe?g|gif|png)/g;
    return url.replace(IMG_URL_REGEX, function(match, href) {
      return (href) ? match : '<img src="' + match + '">';
    });
  };
});

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
