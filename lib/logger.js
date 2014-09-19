var bunyan = require('bunyan');
var logger = bunyan.createLogger({
  name: 'uchiwa',
  src: true
});

module.exports = logger;
