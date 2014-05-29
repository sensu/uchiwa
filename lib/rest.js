var http = require("http");
var https = require("https");
var _ = require("underscore");

function Rest() {
}

Rest.prototype.auth = function(headers, config, callback){
  if(! _.isEmpty(config.user) && ! _.isEmpty(config.pass)){
    var auth = 'Basic ' + new Buffer(config.user + ":" + config.pass).toString('base64');
    _.extend(headers, {'Authorization': auth});
  }
  return headers;
};

Rest.prototype.get = function(options, config, callback) {
  options.headers = this.auth(options.headers, config);
  var prot = options.ssl ? https : http;
  var req = prot.request(options, function(res) {
    var output = "";
    var obj = null;
    res.setEncoding("utf8");
    if (res.statusCode >= 400) {
      callback(res.statusCode, null);
    }
    else {
      res.on("data", function (chunk) {
        output += chunk;
      });
      res.on('end', function() {
        try {
          var jsonResult = JSON.parse(output);
        } catch (e) {
          callback("Error! Data received by Sensu API " + config.host + "was corrupted");
        }
        callback(null, jsonResult);
      });
    }
  });

  req.on("error", function(err) {
    callback(err);
  });
  req.setTimeout(10000, function(){
    req.end();
    req.destroy();
  });

  req.end();
};

Rest.prototype.post = function(options, data, config, callback) {
  options.headers = this.auth(options.headers, config);
  var prot = options.ssl ? https : http;
  var req = prot.request(options, function(res) {
    var output = "";
    var obj = null;
    res.setEncoding("utf8");
    if (res.statusCode >= 400) {
      callback(res.statusCode, null);
    }
    else {
      res.on("data", function (chunk) {
        output += chunk;
      });
      res.on('end', function() {
        try {
          var jsonResult = JSON.parse(output);
        } catch (e) {
          callback("Error! Data received by Sensu API " + config.host + "was corrupted");
        }
        callback(null, jsonResult);
      });
    }
  });

  req.write(data);
  req.end();

  req.on("error", function(err) {
    callback(err);
  });
  req.setTimeout(10000, function(){
    req.end();
    req.destroy();
  });

  req.end();
};

Rest.prototype.delete = function(options, config, callback) {
  options.headers = this.auth(options.headers, config);
  var prot = options.port == 443 ? https : http;
  var req = prot.request(options, function(res) {
    var output = "";
    var obj = null;
    res.setEncoding("utf8");
    if (res.statusCode >= 400) {
      callback(res.statusCode, null);
    }
    else {
      res.read(0);
      res.on('end', function() {
        callback(null);
      });
    }
  });

  req.on("error", function(err) {
    callback(err);
  });
  req.setTimeout(10000, function(){
    req.end();
    req.destroy();
  });

  req.end();
};

exports.Rest = Rest;
