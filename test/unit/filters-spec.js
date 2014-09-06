'use strict';

describe('filters', function () {

  beforeEach(module('uchiwa.filters'));

  describe('encodeURIComponent', function () {

    it('should encode URI', inject(function (encodeURIComponentFilter) {
      expect(encodeURIComponentFilter('dc name/client name?check=check name')).toBe('dc%20name%2Fclient%20name%3Fcheck%3Dcheck%20name');
    }));

  });

  describe('displayObject', function () {

    it('should display object', inject(function (displayObjectFilter) {
      expect(displayObjectFilter('test')).toBe('test');
      expect(displayObjectFilter(['test', 'test1', 'test2'])).toBe('test, test1, test2');
      expect(displayObjectFilter({key: 'value'})).toEqual({key: 'value'});
    }));

  });

  describe('filterSubscriptions', function () {

    it('should filter subscriptions', inject(function (filterSubscriptionsFilter) {
      expect(filterSubscriptionsFilter([
        {name: 'test1', subscriptions: []},
        {name: 'test2', subscriptions: ['linux']}
      ], 'linux')).toEqual([
        {name: 'test2', subscriptions: ['linux']}
      ]);
      expect(filterSubscriptionsFilter([
        {name: 'test1', subscriptions: []},
        {name: 'test2', subscriptions: ['linux']}
      ], '')).toEqual([
        {name: 'test1', subscriptions: []},
        {name: 'test2', subscriptions: ['linux']}
      ]);
    }));

  });
});