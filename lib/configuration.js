'use strict';

var _ = require('underscore');
var yargs = require('yargs')
    .describe('c', 'Load a config (file)')
    .alias('c', 'config')
    .alias('c', 'config-file')
    .alias('c', 'config_file')
    .default('c', './config.json');
var argv = yargs.argv;
var logger = require('./logger.js');

if (argv.c.substring(0, 2) === './') {
  var file =  '.' + argv.c;
}
else {
  var file = argv.c;
}

var configuration = {};

configuration.exist = function (next) {
  try {
    var config = require(file);
    next(null, config);
  } catch (e) {
    next('Could not open ' + file + '. ' + e);
  }
};

configuration.get = function (callback) {
  var self = this;
  this.exist(function(err, result){
    if (err) {
      yargs.showHelp();
      logger.fatal(err);
      process.exit(1);
    }
    else {
      self.initialize(result, function(result){
        callback(result);
      });
    }
  });
};

configuration.initialize = function (config, next) {
  // Backward compatibilty with uchiwa < 0.2.x
  if (!_.isArray(config.sensu)) {
    config.sensu = [config.sensu];
    config.sensu[0].name = config.sensu[0].host;
  }

  // Initialize missing values
  config.uchiwa.port = config.uchiwa.port || 3000;
  config.uchiwa.host = config.uchiwa.host || '0.0.0.0';
  config.uchiwa.refresh = config.uchiwa.refresh || 10000;
  config.uchiwa.dateFormat = config.uchiwa.dateFormat || 'YYYY[-]MM[-]DD HH[:]mm[:]ss';
  config.uchiwa.logLevel = config.uchiwa.logLevel || 'info';

  next(config);
};

configuration.public = function (config) {
  var publicConfig = JSON.parse(JSON.stringify(config));
  publicConfig.uchiwa.user = '*****';
  publicConfig.uchiwa.pass = '*****';
  _.each(publicConfig.sensu, function (element) {
    element.user = '*****';
    element.pass = '*****';
  });
  return publicConfig;
};

module.exports = configuration;
