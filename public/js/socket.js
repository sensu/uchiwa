$(document).ready(function () {
  socket = io.connect();

  socket.on('checks', function(data) {
    checks = JSON.parse(data.content);
  });

  socket.on('stashes', function(data) {
    stashes = JSON.parse(data.content);
  });

  socket.on('messenger', function(data) {
    var message = JSON.parse(data.content);
    notification(message.type, message.content);
  });

  //
  // Clients
  //
  socket.on('clients', function(data) {
    clients = JSON.parse(data.content);
    updateClients(clients);
    updateDashboard();
  });

  //
  // Events
  //
  socket.on('events', function(data) {
    events = JSON.parse(data.content);
  });

  //
  // Client details
  //
  socket.on('client', function(data) {
    client.history = JSON.parse(data.content);
    updateClient(client);
  });
});
