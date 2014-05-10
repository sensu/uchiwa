$(document).ready(function () {
  socket = io.connect();

  socket.on('checks', function(data) {
    if(_.isUndefined(data.content)) return;
    checks = JSON.parse(data.content);
    updateChecks(checks);
  });

  socket.on('stashes', function(data) {
    if(_.isUndefined(data.content)) return;
    stashes = JSON.parse(data.content);
    updateStashes(stashes);
  });

  socket.on('messenger', function(data) {
    if(_.isUndefined(data.content)) return;
    var message = JSON.parse(data.content);
    notification(message.type, message.content);
  });

  //
  // Clients
  //
  socket.on('clients', function(data) {
    if(_.isUndefined(data.content)) return;
    clients = JSON.parse(data.content);
    updateClients(clients);
    updateDashboard();
  });

  //
  // Events
  //
  socket.on('events', function(data) {
    if(_.isUndefined(data.content)) return;
    events = JSON.parse(data.content);
    updateEvents(events);
  });

  //
  // Client details
  //
  socket.on('client', function(data) {
    if(_.isUndefined(data.content)) return;
    client.history = JSON.parse(data.content);
    updateClient(client);
  });
});
