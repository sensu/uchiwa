function Check(data) {
  this.name = data.name;
  this.command = data.command;
  this.subscribers = data.subscribers;
  this.interval = data.interval;
}

Check.prototype.isSilenced = function(path, callback){
  var result = stashes.filter(function (e) { return e.path === path });
  if(result.length > 0){
    callback(true);
  }
  else {
    callback(null);
  }
}