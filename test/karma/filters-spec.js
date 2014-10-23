'use strict';

describe('filters', function () {

  var $filter;
  var settings;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function (_socket_, _$filter_, _settings_) {
    $filter = _$filter_;
    settings = _settings_;
  }));

  describe('buildStashes', function () {

    it('should add client & check properties', inject(function (buildStashesFilter) {
      var stashes = [
        {path: 'silence/foo/bar'},
        {path: 'silence/baz'}
      ];
      var expectedStashes = [
        {client: 'foo', check: 'bar', path: 'silence/foo/bar'},
        {client: 'baz', check: null, path: 'silence/baz'}
      ];
      expect(buildStashesFilter(stashes)).toEqual(expectedStashes);
      expect(buildStashesFilter('string')).toEqual('string');
      expect(buildStashesFilter({})).toEqual({});
    }));

  });

  describe('displayObject', function () {

    it('should display object', inject(function (displayObjectFilter) {
      expect(displayObjectFilter('test')).toBe('test');
      expect(displayObjectFilter(['test', 'test1', 'test2'])).toBe('test, test1, test2');
      expect(displayObjectFilter({key: 'value'})).toEqual({key: 'value'});
    }));

  });

  describe('getTimestamp', function () {

    it('should convert epoch to human readable date', inject(function (getTimestampFilter) {
      expect(getTimestampFilter('test')).toBe('test');
      expect(getTimestampFilter(1)).toBe(1);
      expect(getTimestampFilter(1410908218)).toBe(moment.utc('2014-09-16 22:56:58', 'YYYY-MM-DD HH:mm:ss').local().format('YYYY-MM-DD HH:mm:ss'));
    }));

  });

  describe('getExpireTimestamp', function () {

    it('should convert epoch to human readable date', inject(function (getExpireTimestampFilter) {
      expect(getExpireTimestampFilter('test')).toBe('Unknown');
      expect(getExpireTimestampFilter(900, 1410908218)).toBe(moment.utc('2014-09-16 23:11:58').local().format('YYYY-MM-DD HH:mm:ss'));
      expect(getExpireTimestampFilter(-1, 1410908218)).toBe('Never');
    }));

  });

  describe('encodeURIComponent', function () {

    it('should encode URI', inject(function (encodeURIComponentFilter) {
      expect(encodeURIComponentFilter('dc name/client name?check=check name')).toBe('dc%20name%2Fclient%20name%3Fcheck%3Dcheck%20name');
    }));

  });

  describe('filterSubscriptions', function () {

    it('should filter subscriptions', inject(function (filterSubscriptionsFilter, $filter, settings) {
      expect(filterSubscriptionsFilter([{name: 'test1', subscriptions: []}, {name: 'test2', subscriptions: ['linux']}], 'linux')).toEqual([{name: 'test2', subscriptions: ['linux']}]);
      expect(filterSubscriptionsFilter([{name: 'test1', subscriptions: []}, {name: 'test2', subscriptions: ['linux']}], '')).toEqual([{name: 'test1', subscriptions: []}, {name: 'test2', subscriptions: ['linux']}]);
    }));

  });

  describe('getStatusClass', function () {

    it('should return CSS class based on status', inject(function (getStatusClassFilter) {
      expect(getStatusClassFilter(0)).toBe('success');
      expect(getStatusClassFilter(1)).toBe('warning');
      expect(getStatusClassFilter(2)).toBe('critical');
      expect(getStatusClassFilter(3)).toBe('unknown');
      expect(getStatusClassFilter('foo')).toBe('unknown');
    }));

  });

  describe('getAckClass', function () {

    it('should return icon based on acknowledgment', inject(function (getAckClassFilter) {
      expect(getAckClassFilter(true)).toBe('fa-volume-off');
      expect(getAckClassFilter(null)).toBe('fa-volume-up');
    }));

  });

});
