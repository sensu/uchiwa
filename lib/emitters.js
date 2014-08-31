'use strict';

var bunyan = require('bunyan');
var log = bunyan.createLogger({name: 'uchiwa', src: true});

var emitters = {};

var messages = {
  createStash: {
    error: '<strong>Error!</strong> The stash was not created: ',
    success: '<strong>Success!</strong> The stash has been created.'
  },
  deleteClient: {
    error: '<strong>Error!</strong> The client was not deleted: ',
    success: '<strong>Success!</strong> The client has been deleted.'
  },
  deleteStash: {
    error: '<strong>Error!</strong> The stash was not deleted: ',
    success: '<strong>Success!</strong> The stash has been deleted.'
  },
  generic: {
    error: '<strong>Error!</strong> ',
    success: '<strong>Success!</strong>'
  },
  resolveEvent: {
    error: '<strong>Error!</strong> The event was not resolved: ',
    success: '<strong>Success!</strong> The event has been resolved.'
  }
};

emitters.alert = function(req, err, object) {
  var type = (err) ? 'error' : 'success';
  var message = messages[object][type];
  if (err) { message += err; }

  if (req) {
    req.io.emit('messenger', {
      content: JSON.stringify({
        'type': type,
        'content': message
      })
    });
  }
  else {
    log.warn('Could not emit to socket client "'+ req.io.id +'"');
  }

};

emitters.send = function (req, err, result, object) {
  if (err) {
    this.alert(req, err, 'generic');
  }
  else {
    if (req) {
      req.io.emit(object, {content: JSON.stringify(result)});
    }
    else {
      log.warn('Could not emit to socket client "'+ req.io.socket.id +'"');
    }
  }
};

module.exports = emitters;