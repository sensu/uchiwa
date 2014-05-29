var moment = require('moment');

function Stats(config, sensu) {
  this.config = config;
  this.sensu = sensu;
  this.dashboard = [];
}

Stats.prototype.getDashboard = function(callback){
  var retention = this.config.retention || 10;
  this.dashboard.unshift({y: moment().format('YYYY[-]MM[-]DD HH[:]mm[:]ss'), e: this.sensu.events.length, s: this.sensu.stashes.length});
  if(this.dashboard.length > (retention * 6)) this.dashboard.pop();
  callback(null);
};

exports.Stats = Stats;