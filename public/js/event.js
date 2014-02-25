function Event(data) {
  this.check = data.check;
  this.client = data.client;
  this.flapping = data.flapping;
  this.handlers = data.handlers;
  this.issued = data.issued;
  this.last_check = data.last_check;
  this.occurrences = data.occurrences;
  this.output = data.output;
  this.status = data.status;
}