function Client(id) {
  this.name = id;
  this.history = [];
}

Client.prototype.getEvents = function() {
  var name = this.name;
  return events.filter(function (e) { return e.client == name });
}

Client.prototype.getSubscription = function(callback) {
  var name = this.name;
  var subscription = clients.filter(function (e) { return e.name == name });
  if(subscription.length > 0){
    callback(null, subscription);
  }
  else {
    callback(true);
  }
}