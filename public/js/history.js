function History(data) {
  this.check = data.check;
  this.history = data.history;
  this.last_execution = data.last_execution;
  this.last_status = data.last_status;
  this.last_check = data.last_check;
}

History.prototype.getEvent = function(events, client, callback) {
  if(this.last_status == 0){
    callback(null, null);
  }
  else {
    var check = this.check;
    var details = events.filter(function (e) { return e.client == client && e.check == check });
    if (details.length == 1){
      event = new Event(details[0]);
      callback(null, event);
    }
    else {
      callback(true, "");
    }
  }
}

History.prototype.getStyle = function (callback) {
  if (this.last_status == 2){
    callback("danger");
  }
  else if (this.last_status == 1){
    callback("warning");
  }
  else {
    callback("success");
  }
}

History.prototype.getCheck = function (callback) {
   checkName = this.check;
  if(checkName == "keepalive"){
    callback(null, new Check({name: checkName, command: "keepalive", subscribers: "", interval: ""}));
  }
  else {
    //console.log(checkName + " " + checks);
    var check = checks.filter(function (e) { return e.name == checkName });
    
    if (check.length > 0){
      callback(null, new Check(check[0]));
    }
    else {
      callback(true, "");
    }
  }
  
} 