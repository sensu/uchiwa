'use strict';

var moment = require('moment');

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

emitters.alert = function(socket, err, object) {
  var type = (err) ? 'error' : 'success';
  var message = messages[object][type];
  if (err) { message += err; }

  if (socket) {
    socket.emit('messenger', {
      content: JSON.stringify({
        'type': type,
        'content': message
      })
    });
  }
  else {
    console.log(moment().format() + ' [info] Could not emit to socket client "'+ socket.id +'"');
  }
  
};

emitters.send = function (socket, err, result, object) {
  if (err) {
    this.alert(socket, err, 'generic');
  }
  else {
    if (socket) {
      socket.emit(object, {content: JSON.stringify(result)});
    }
    else {
      console.log(moment().format() + ' [info] Could not emit to socket client "'+ socket.id +'"');
    } 
  }
};

module.exports = emitters;