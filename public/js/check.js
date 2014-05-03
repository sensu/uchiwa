function Check(data) {
  this.name = data.name;
  this.command = data.command;
  this.subscribers = (_.isUndefined(data.subscribers)) ? "<em>n/a</em>" : data.subscribers;
  this.interval = data.interval;
  this.type = (_.isUndefined(data.type)) ? "normal" : data.type;
  this.handlers = (_.isUndefined(data.handlers)) ? "<em>n/a</em>" : data.handlers;
  this.handle = (_.isUndefined(data.handle)) ? "true" : data.handle;
  this.subdue = (_.isUndefined(data.subdue)) ? {"begin": "<em>n/a</em>", "end": "<em>n/a</em>"} : data.subdue;
  this.standalone = (_.isUndefined(data.standalone)) ? "false" : data.standalone;
}

Check.prototype.isSilenced = function(path){
  var result = stashes.filter(function (e) { return e.path === path });
  if(result.length > 0){
    return true;
  }
  else {
    return false;
  }
}