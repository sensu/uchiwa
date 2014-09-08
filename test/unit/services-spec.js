'use strict';

describe('services', function () {
  var socket;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function (_socket_) {
    socket = _socket_;
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

  describe('stashesService', function () {

    it('should have a stash method', inject(function (stashesService) {
      expect(stashesService.stash).toBeDefined();
    }));

    describe('stash()', function () {

      it('should emit delete_stash', inject(function (stashesService) {
        var mockPayload = {dc: 'dcName', payload: {path: '/', content: {}}};
        spyOn(socket, 'emit');
        stashesService.stash(mockPayload.dc, mockPayload.payload);
        expect(socket.emit).toHaveBeenCalledWith('delete_stash', JSON.stringify(mockPayload));
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