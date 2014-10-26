'use strict';

describe('filters', function () {

  var $filter;
  var settings;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function (_$filter_, _settings_) {
    $filter = _$filter_;
    settings = _settings_;
  }));

  describe('arrayLength', function () {

    it('should return 0 if null', inject(function (arrayLengthFilter) {
      expect(arrayLengthFilter(null)).toEqual(0);
    }));

    it('should return 0 if not an array', inject(function (arrayLengthFilter) {
      expect(arrayLengthFilter('string')).toEqual(0);
    }));

    it('should return proper array length', inject(function (arrayLengthFilter) {
      expect(arrayLengthFilter([0, 1, 2])).toEqual(3);
    }));

  });

  describe('arrayToString', function () {

    it('should return 0 if null', inject(function (arrayToStringFilter) {
      expect(arrayToStringFilter(null)).toEqual('');
    }));

    it('should return 0 if not an array', inject(function (arrayToStringFilter) {
      expect(arrayToStringFilter('string')).toEqual('string');
    }));

    it('should return proper array length', inject(function (arrayToStringFilter) {
      expect(arrayToStringFilter([0, 1, 2])).toEqual('0 1 2');
    }));

  });

  describe('buildEvents', function () {

    it('should not accept anything else than an array', inject(function (buildEventsFilter) {
      expect(buildEventsFilter('string')).toEqual('string');
      expect(buildEventsFilter({})).toEqual({});
    }));

    it('handles missing check OR client objects', inject(function (buildEventsFilter) {
      var events = [
        { check: { source: 'foo' }},
        { client: {name: 'baz' }}
      ];
      var expectedEvents = [
        { check: { source: 'foo'}, sourceName: 'foo'},
        { check: {}, client: {name: 'baz'}, sourceName: 'baz'}
      ];
      expect(buildEventsFilter(events)).toEqual(expectedEvents);
    }));

    it('handles missing check AND client objects', inject(function (buildEventsFilter) {
      var events = [{}];
      var expectedEvents = [{sourceName: 'unknown'}];
      expect(buildEventsFilter(events)).toEqual(expectedEvents);
    }));

    it('should add sourceName properties', inject(function (buildEventsFilter) {
      var events = [
        { check: { source: 'foo'}, client: { name: 'bar'}},
        { check: { name: 'qux'}, client: {name: 'baz'}}
      ];
      var expectedEvents = [
        { check: { source: 'foo'}, client: { name: 'bar'}, sourceName: 'foo'},
        { check: { name: 'qux'}, client: {name: 'baz'}, sourceName: 'baz'}
      ];
      expect(buildEventsFilter(events)).toEqual(expectedEvents);
    }));

  });

  describe('buildStashes', function () {

    it('should not accept anything else than an array', inject(function (buildStashesFilter) {
      expect(buildStashesFilter('string')).toEqual('string');
      expect(buildStashesFilter({})).toEqual({});
    }));

    it('should add client & check properties', inject(function (buildStashesFilter) {
      var stashes = [
        {path: 'silence/foo/bar'},
        {path: 'silence/'},
        {path: 'silence/baz'}
      ];
      var expectedStashes = [
        {client: 'foo', check: 'bar', path: 'silence/foo/bar'},
        {client: null, check: null, path: 'silence/'},
        {client: 'baz', check: null, path: 'silence/baz'}
      ];
      expect(buildStashesFilter(stashes)).toEqual(expectedStashes);
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

  describe('richOutput', function () {

    it('should convert an object to JSON string', inject(function (richOutputFilter) {
      expect(richOutputFilter({foo: 'bar'})).toEqual('{&#34;foo&#34;:&#34;bar&#34;}');
    }));

    it('should convert an image URL to a HTML image', inject(function (richOutputFilter) {
      expect(richOutputFilter('http://foo.bar/baz.gif')).toContain('<a target="_blank" href="http://foo.bar/baz.gif"><img src=');
    }));

    it('should convert an URL to a HTML URL', inject(function (richOutputFilter) {
      expect(richOutputFilter('http://foo.bar/baz')).toContain('<a target="_blank" href="http://foo.bar/baz">');
    }));

  });

  describe('setMissingProperty', function () {

    it('should set to false a missing property', inject(function (setMissingPropertyFilter) {
      expect(setMissingPropertyFilter(undefined)).toBe(false);
      expect(setMissingPropertyFilter({foo: 'bar'})).toEqual({foo: 'bar'});
    }));

  });

});
