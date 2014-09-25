'use strict';

var expect = require('chai').expect;

var health = require('../../lib/health.js');

var config = {};
var sensu = {};

beforeEach(function(){
  config = {};
  sensu = {};
})

describe('health.sensu', function () {
  it('should return HTTP 200 with proper payload when everything is okay', function (done) {
    config.sensu = [{ name: 'foo' }];
    sensu.dc = [{ name: 'foo', health: 'ok' }];

    health.sensu(sensu, config, function (result) {
      expect(result).to.eql({ json: {"foo": {"output": "ok" }}, code: 200 });
      done();
    });

  });

  it('should return HTTP 404 when config object is invalid', function (done) {
    sensu.dc = [{ name: 'foo', health: 'ok' }];
    health.sensu(sensu, config, function (result) {
      expect(result).to.eql({ json: {}, code: 404 });
      done();
    });
  });

  it('should return HTTP 404 when sensu object is invalid', function (done) {
    config.sensu = [{ name: 'foo' }];
    health.sensu(sensu, config, function (result) {
      expect(result).to.eql({ json: {}, code: 404 });
      done();
    });
  });

  it('should return a payload that indicate missing datacenters', function (done) {
    config.sensu = [{ name: 'foo' }, { name: 'bar' }];
    sensu.dc = [{ name: 'foo', health: 'ok' }];
    var expectedPayload = {
      json:
        {
          foo: { output: 'ok' },
          bar: { output: 'error: bar is missing' }
        },
        code: 200
    };

    health.sensu(sensu, config, function (result) {
      expect(result).to.eql(expectedPayload);
      done();
    });
  });
});
