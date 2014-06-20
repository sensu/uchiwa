'use strict';

var moment = require('moment');
var _ = require('underscore');

function Stats(config) {
  this.config = config;
  this.dashboard = [];
}

Stats.prototype.getDashboard = function (sensu) {
  var retention = this.config.retention || 10;
  var retentionPoints = retention * 60000 / this.config.refresh;
  var stashes = 0;
  var events = 0;
  _.each(sensu.events, function (element) {
    events += element.length;
  });
  _.each(sensu.stashes, function (element) {
    stashes += element.length;
  });

  this.dashboard.unshift({y: moment().format('YYYY[-]MM[-]DD HH[:]mm[:]ss'), e: events, s: stashes});
  if (this.dashboard.length > retentionPoints) {
    this.dashboard.pop();
  }
};

exports.Stats = Stats;