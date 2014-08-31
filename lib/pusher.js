'use strict';

var async = require('async');
var _ = require('underscore');

var pusher = {};

pusher.pull = function (app, sensu, datacenters, callback) {
  var attributes = ['checks', 'clients', 'events', 'stashes'];
  attributes.forEach(function(attribute) {
    sensu[attribute] = [];
  });
  sensu.dc = [];

  async.eachSeries(datacenters, function (datacenter, nextDc) {
    datacenter.pull(function () {
      var aggregate = function (callback) {
        async.each(attributes, function (attribute, nextAttribute) {
          async.each(datacenter.sensu[attribute], function (item, nextItem) {
            item.dc = datacenter.name;
            sensu[attribute].push(item);
            nextItem();
          }, function () {
            nextAttribute();
          });
        }, function () {
          callback();
        });
      };
      aggregate(function () {
        datacenter.build();
        sensu.dc.push({
          name: datacenter.name,
          style: datacenter.style,
          clients: datacenter.clients,
          events: datacenter.events,
          stashes: datacenter.stashes,
          checks: datacenter.checks,
          info: datacenter.info,
          health: datacenter.health
        });
        nextDc();
      });
    },
    function (messageContent) {
      app.io.broadcast('messenger', {
        content: messageContent
      });
    }
  );
  }, function () {
    sensu.subscriptions = [];
    async.each(sensu.clients, function (client, nextClient) {
      if(_.isObject(client.subscriptions)){
        async.each(client.subscriptions, function (subscription, nextSubscription) {
          if(sensu.subscriptions.indexOf(subscription) === -1) { sensu.subscriptions.push(subscription); }
          nextSubscription();
        });
      }
      nextClient();
    }, function() {
      callback(sensu);
    });
  });
};

pusher.push = function (app, sensu) {
  app.io.broadcast('sensu', {content: JSON.stringify(sensu)});
};

module.exports = pusher;