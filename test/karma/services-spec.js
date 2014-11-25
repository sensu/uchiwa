'use strict';

describe('services', function () {
  var $rootScope;
  var $scope;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function (_$rootScope_) {
    $rootScope = _$rootScope_;
    $scope = $rootScope.$new();
  }));

  describe('Page', function () {

    it('should have a title method', inject(function (Page) {
      expect(Page.title).toBeDefined();
    }));

    describe('title()', function () {
      it('should suffix the application title', inject(function (Page) {
        var title = 'Test';
        Page.setTitle(title);
        expect(Page.title()).toBe(title + ' | Uchiwa');
      }));
    });

  });

  describe('clientsService', function () {

    describe('resolveEvent', function () {

      it('should return false when neither client & check are undefined or not objects', inject(function (clientsService) {
        expect(clientsService.resolveEvent('foo', null, null)).toEqual(false);
        expect(clientsService.resolveEvent('foo', undefined)).toEqual(false);
      }));

      it('should emit HTTP POST to /resolveEvent', inject(function (clientsService, uchiwaBackend) {
        var expectedPayload = {dc: 'foo', payload: {client: 'bar', check: 'qux'}};
        spyOn(uchiwaBackend, 'resolveEvent').and.callThrough();
        clientsService.resolveEvent('foo', {name: 'bar'}, {check: 'qux'});
        expect(uchiwaBackend.resolveEvent).toHaveBeenCalledWith(expectedPayload);
      }));

    });

  });

  describe('stashesService', function () {

    it('should have a deleteStash method', inject(function (stashesService) {
      expect(stashesService.deleteStash).toBeDefined();
    }));

    describe('submit', function () {

      it('should emit HTTP POST to /createStash for a client', inject(function (stashesService, uchiwaBackend) {
        spyOn(uchiwaBackend, 'createStash').and.callThrough();
        var timestamp = Math.floor(new Date()/1000);
        var expectedPayload = {dc: 'foo', payload: {path: 'silence/bar', content: {reason: '', source: 'uchiwa', timestamp: timestamp}}};
        stashesService.submit({name: 'bar', acknowledged: false, dc: 'foo'}, {path: ['bar', '']});
        expect(uchiwaBackend.createStash).toHaveBeenCalledWith(expectedPayload);
      }));

      it('should emit HTTP POST to /createStash for a check', inject(function (stashesService, uchiwaBackend) {
        spyOn(uchiwaBackend, 'createStash').and.callThrough();
        var timestamp = Math.floor(new Date()/1000);
        var expectedPayload = {dc: 'foo', payload: {path: 'silence/bar/qux', content: {reason: '', source: 'uchiwa', timestamp: timestamp}}};
        stashesService.submit({client: {name: 'bar'}, check: {name: 'qux'}, acknowledged: false, dc: 'foo'}, {path: ['bar', 'qux']});
        expect(uchiwaBackend.createStash).toHaveBeenCalledWith(expectedPayload);
      }));

    });

  });

  describe('routingService', function () {

    it('should have a go method', inject(function (routingService) {
      expect(routingService.go).toBeDefined();
    }));
    it('should have a deleteEmptyParameter method', inject(function (routingService) {
      expect(routingService.deleteEmptyParameter).toBeDefined();
    }));
    it('should have a initFilters method', inject(function (routingService) {
      expect(routingService.initFilters).toBeDefined();
    }));
    it('should have a permalink method', inject(function (routingService) {
      expect(routingService.permalink).toBeDefined();
    }));
    it('should have a updateFilters method', inject(function (routingService) {
      expect(routingService.updateFilters).toBeDefined();
    }));
    it('should have a updateValue method', inject(function (routingService) {
      expect(routingService.updateValue).toBeDefined();
    }));

    describe('go()', function() {

      it('should call $location.url', inject(function (routingService, $location) {
        var uri = '/testing';
        spyOn($location, 'url');
        routingService.go(uri);
        expect($location.url).toHaveBeenCalledWith(uri);
      }));

      it('should encode URIs', inject(function (routingService, $location) {
        var uri = '/this needs !@#$ encoding';
        spyOn($location, 'url');
        routingService.go(uri);
        expect($location.url).not.toHaveBeenCalledWith(uri);
        expect($location.url).toHaveBeenCalledWith(encodeURI(uri));
      }));

    });

  });

  describe('underscore', function () {

    it('should define _', inject(function (underscore) {
      expect(underscore).toBe(window._);
    }));

  });

});
