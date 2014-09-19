'use strict';

var expect = require('chai').expect;

var Sensu = require('../../lib/sensu.js').Sensu
var config = {
  host: '127.0.0.1'
};
var sensu = new Sensu(config);

beforeEach(function(){
  sensu.checks = [];
  sensu.clients = [];
  sensu.events = [];
  sensu.stashes = [];
})

describe('findStash', function () {
  it('should return stash with client & check name', function (done) {
    sensu.stashes = [{ path: 'silence/foo/check_bar' }];
    var result = sensu.findStash('foo', 'check_bar');
    expect(result).to.equal(true);
    done();
  });

  it('should return stash client name', function (done) {
    sensu.stashes = [{ path: 'silence/foo' }];
    var result = sensu.findStash('foo');
    expect(result).to.equal(true);
    done();
  });
});

describe('buildClients', function () {
  it('should return a new client object', function (done) {
    var history = [
      {
        check: 'check_bar',
        history: [0, 1],
        last_execution: 1370725352,
        last_status: 1
      },
      {
        check: 'check_baz',
        history: [1, 0],
        last_execution: 1370725353,
        last_status: 2
      }
    ];
    sensu.checks = [
      {
        name: 'check_bar',
        command: 'check_bar.sh'
      },
      {
        name: 'check_baz',
        command: 'check_baz.sh'
      }
    ];
    sensu.clients = [
      { name: 'foo' },
      { name: 'bar' }
    ];
    sensu.events = [
      {
        client: {
          name: 'foo'
        },
        check: {
          name: 'check_baz',
          output: 'output of baz'
        }
      },
    ];
    sensu.stashes = [
      { "path": "silence/foo/check_baz" }
    ];
    var expectedClient = {
      name: 'foo',
      history: [
        {
          check: 'check_bar',
          history: [0, 1],
          last_execution: 1370725352,
          last_status: 1,
          acknowledged: null,
          output: '',
          model: {
            name: 'check_bar',
            command: 'check_bar.sh'
          }
        },
        {
          check: 'check_baz',
          history: [1, 0],
          last_execution: 1370725353,
          last_status: 2,
          acknowledged: true,
          output: 'output of baz',
          model: {
            name: 'check_baz',
            command: 'check_baz.sh'
          }
        }
      ]
    };

    sensu.buildClient('foo', history, function(err, result){
      expect(result).to.eql(expectedClient);
      done();
    });
  });
});

describe('buildClients', function () {
  it('should add acknowledged, eventsSummary, status & version properties', function (done) {
    sensu.clients = [
      { name: 'foo' },
      { name: 'bar' },
      { name: 'baz' },
      { name: 'qux', version: '0.13.1' }
    ];
    sensu.events = [
      { client: {name: 'foo'}, check: {name: 'foo_event', status: 1} },
      { client: {name: 'bar'}, check: {name: 'bar_warning_event', status: 1} },
      { client: {name: 'bar'}, check: {name: 'bar_critical_event', status: 2} },
      { client: {name: 'baz'}, check: {name: 'bar_unknown_event', status: 3} }
    ];
    sensu.stashes = [
      { "path": "silence/baz" }
    ];
    var expectedClients = [
      { name: 'foo', version: '0.12.x', status: 1, eventsSummary: 'foo_event', acknowledged: null },
      { name: 'bar', version: '0.12.x', status: 2, eventsSummary: 'bar_warning_event and 1 more...', acknowledged: null },
      { name: 'baz', version: '0.12.x', status: 3, eventsSummary: 'bar_unknown_event', acknowledged: true },
      { name: 'qux', version: '0.13.1', status: 0, eventsSummary: '', acknowledged: null }
    ];

    sensu.buildClients(function(){
      expect(sensu.clients).to.eql(expectedClients);
      done();
    });
  });
});

describe('buildEvents', function () {
  it('should convert a Sensu 0.12.x event object', function (done) {
    sensu.events = [
      { client: 'foo', check: 'foo_event', issued: 123456789, status: 1 },
      { client: 'bar', check: 'bar_event', flapping: true, occurrences: 10, output: 'bar_event' }
    ];
    sensu.stashes = [
      { "path": "silence/bar/bar_event" }
    ];
    var expectedEvents = [
      { client: { name: 'foo' }, check: { name: 'foo_event', issued: 123456789, status: 1 }, occurrences: 1, action: 'create', acknowledged: null },
      { client: { name: 'bar' }, check: { name: 'bar_event', output: 'bar_event' }, occurrences: 10, action: 'flapping', acknowledged: true }
    ];

    sensu.buildEvents(function(){
      expect(sensu.events).to.eql(expectedEvents);
      done();
    });
  });

  it('should only add acknowledged property to a Sensu >= 0.13.x event object', function (done) {
    sensu.events = [
      { id: 'abc', client: { name: 'foo' }, check: { name: 'foo_event', issued: 123456789, status: 1 }, occurrences: 1, action: 'create'},
    ];
    sensu.stashes = [
      { "path": "silence/foo/foo_event" }
    ];
    var expectedEvents = [
      { id: 'abc', client: { name: 'foo' }, check: { name: 'foo_event', issued: 123456789, status: 1 }, occurrences: 1, action: 'create', acknowledged: true }
    ];

    sensu.buildEvents(function(){
      expect(sensu.events).to.eql(expectedEvents);
      done();
    });
  });
});
