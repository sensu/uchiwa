'use strict';

var basicAuth = require('basic-auth');

var authentication = {};

authentication.basic = function (req, res, next) {
  var config = res.app.set('config');
  var user = basicAuth(req);

  function unauthorized(res) {
    res.set('WWW-Authenticate', 'Basic realm=Authorization Required');
    return res.send(401);
  }
  if (req.path.substring(0, 7) === '/health') {
    next();
  }
  else {
    if (!user || !user.name || !user.pass) {
      return unauthorized(res);
    }
    if (user.name === config.uchiwa.user && user.pass === config.uchiwa.pass) {
      next();
    }
    else {
      return unauthorized(res);
    }
  }
};

module.exports = authentication;