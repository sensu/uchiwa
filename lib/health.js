'use strict';

var _ = require('underscore');
var logger = require('./logger.js');

var health = {};

health.get = function (req, res, sensu, config) {
  var self = this;
  this.sensu(sensu, config, function (result) {
    if (req.params.component) {
      if (req.params.component === 'uchiwa') {
        self.respond(res, {uchiwa: 'ok'}, 200);
      }
      else if (req.params.component === 'sensu') {
        self.respond(res, {sensu: result.json}, result.code);
      }
      else {
        self.respond(res, {component: 'not found'}, 404);
      }
    }
    else {
      self.respond(res, {uchiwa: 'ok', sensu: result.json}, result.code);
    }
  });
};

health.respond = function (res, data, code) {
  res.setHeader('Last-Modified', (new Date()).toUTCString());
  res.setHeader('Content-Type', 'application/json');
  res.send(code, JSON.stringify(data));
};

health.sensu = function (sensu, config, next) {
  var result = {json: {}, code: 200};
  if (!_.isArray(config.sensu) || config.sensu.length === 0) {
    logger.error('Could not retrieve Sensu APIs details from config.json');
    return next({json: {}, code: 404});
  }
  if (!_.isArray(sensu.dc) || sensu.dc.length === 0) {
    logger.error('Datacenters details and statistics are missing');
    return next({json: {}, code: 404});
  }
  config.sensu.forEach(function(api){
    var datacenter = sensu.dc.filter(function (e){
      return e.name === api.name;
    })[0];

    if (datacenter === undefined) {
      result.json[api.name] = {};
      result.json[api.name].output = 'error: ' + api.name + ' is missing';
    }
    else {
      result.json[datacenter.name] = {};
      if (datacenter.health !== 'ok') {
        result.code = 503;
      }
      result.json[datacenter.name].output = datacenter.health;
    }
  });
  next(result);
};

module.exports = health;
