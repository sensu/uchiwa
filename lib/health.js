'use strict';

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
  config.sensu.forEach(function(api){
    var datacenter = sensu.dc.filter(function (e){
      return e.name === api.name;
    })[0];
    result.json[datacenter.name] = {};
    if (datacenter.health !== 'ok') {
      result.code = 503;
    }
    result.json[datacenter.name].output = datacenter.health;
  });
  next(result);
};

module.exports = health;