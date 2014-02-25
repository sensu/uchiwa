var http = require("http");
var https = require("https");

function Rest() {
}

Rest.prototype.get = function(options, callback) {
  var prot = options.port == 443 ? https : http;
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
        var jsonResult = JSON.parse(output);
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

exports.Rest = Rest;