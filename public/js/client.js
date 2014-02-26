function Client(data) {
  this.name = data.name;
  this.address = data.address;
  this.subscriptions = data.subscriptions;
  this.timestamp = data.timestamp;
  this.last_check = data.last_check;
  this.status = data.status;
  this.events = data.events;
  this.history = [];
}

Client.prototype.getEvents = function(callback){
  var name = this.name;
  var clientEvents = events.filter(function (e) { return e.client == name });
  if(clientEvents.length > 0){
    callback(null, clientEvents);
  }
  else {
    callback(true);
  }
}

Client.prototype.getSubscription = function(callback){
  var name = this.name;
  var subscription = clients.filter(function (e) { return e.name == name });
  if(subscription.length > 0){
    callback(null, subscription);
  }
  else {
    callback(true);
  }
}

Client.prototype.eventsCount = function(){
  if (typeof this.events === 'undefined'){
    return "";
  }
  else {
    if (this.events.length != 1){
      return this.events[0].check + " and " + (this.events.length - 1) + " more...";
    }
    else {
      return this.events[0].check;
    }
  }
}