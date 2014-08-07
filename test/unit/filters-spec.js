'use strict';

describe('filter', function() {

  beforeEach(module('uchiwa.filters'));

  describe('encodeURIComponent', function() {

    it('should encore URI', inject(function(encodeURIComponentFilter) {
      expect(encodeURIComponentFilter('dc name/client name?check=check name')).toBe('dc%20name%2Fclient%20name%3Fcheck%3Dcheck%20name');
    }));

  });

  describe('displayObject', function() {

    it('should encore URI', inject(function(displayObjectFilter) {
      expect(displayObjectFilter('test')).toBe('test');
      expect(displayObjectFilter(['test', 'test1', 'test2'])).toBe('test, test1, test2');
      expect(displayObjectFilter({key: 'value'})).toEqual({key: 'value'});
    }));
    
  });
});