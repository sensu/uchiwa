function Stash(data) {
  this.path = data.path;
  this.content = data.content;
  this.expire = data.expire;
  this.last_check = (_.isUndefined(data.last_check)) ? null : data.last_check;

  var path = this.path.split('/');
  this.client = (_.isUndefined(path[1])) ? null : path[1];
  this.check = (_.isUndefined(path[2])) ? null : path[2];
}
