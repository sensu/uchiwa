'use strict';

function Messenger() {
}

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

Messenger.prototype.alert = function(socket, err, object) {
  var type = (err) ? 'error' : 'success';
  var message = messages[object][type];
  if (err) { message += err; }

  socket.emit('messenger', {
    content: JSON.stringify({
      'type': type,
      'content': message
    })
  });
};

Messenger.prototype.post = function (socket, err, result, object) {
  if (err) {
    this.alert(socket, err, 'generic');
  }
  else {
    socket.emit(object, {content: JSON.stringify(result)});
  }
};

exports.Messenger = Messenger;