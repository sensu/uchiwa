'use strict';

var express = require('express');

var authentication = {};

authentication.basic = function (req, res, next) {
  var config = res.app.set('config');
  var basicAuth = express.basicAuth(function (username, password) {
    return (username === config.uchiwa.user && password === config.uchiwa.pass);
  }, 'Restricted area, please identify');

  if (req.path.substring(0, 7) === '/health') { 
    next();
  }
  else {
    return basicAuth.apply(this, arguments);
  }
};

module.exports = authentication;