'use strict';

var logger = require('./logger.js');

var emitters = {};

var messages = {
  createStash: {
    error: 'The stash was not created: ',
    success: 'The stash has been created.'
  },
  deleteClient: {
    error: 'The client was not deleted: ',
    success: 'The client has been deleted.'
  },
  deleteStash: {
    error: 'The stash was not deleted: ',
    success: 'The stash has been deleted.'
  },
  generic: {
    error: 'Error!',
    success: 'Success!'
  },
  resolveEvent: {
    error: 'The event was not resolved: ',
    success: 'The event has been resolved.'
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
    logger.warn('Could not emit to socket client "'+ req.io.id +'"');
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
      logger.warn('Could not emit to socket client "'+ req.io.socket.id +'"');
    }
  }
};

module.exports = emitters;
