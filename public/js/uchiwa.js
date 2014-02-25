var countEvents = function(events, callback) {
  if (typeof events === 'undefined'){
    callback("");
  }
  else {
    if (events.length != 1){
      callback(events[0].check + " and " + (events.length - 1) + " more...");
    }
    else {
      callback(events[0].check);
    }
  }
}

var findEvent = function(events_list, check, client, callback){
  var eventDetails = events_list.filter(function (e) { return e.client == client && e.check == check });
  if (eventDetails.length != 0){
    callback(null, eventDetails);
  }
  else {
    callback(true, "");
  }
  
};

var getStyle = function(status, callback){
  if (status == 2){
    callback("danger");
  }
  else if (status == 1){
    callback("warning");
  }
  else {
    callback("success");
  }
};

/**
* getClient: Request to socket the client detail while modal window is shown
* @param id {string} Name of the client
*/
var getClient = function(id){
  clearInterval(timer);
  client = new Client(id);
  socket.emit('get_client', {name: id});
  // Listen to hide event of modal
  $('#client-details').on('hide.bs.modal', function () {
    $(this).off('hide.bs.modal');
    clearInterval(timer);
  })
  // Fetch client while modal is shown
  var timer = setInterval(function(){
    if($("#client-details").data('bs.modal').isShown){
      socket.emit('get_client', {name: id});
    }
  },  10000);
};

$(document).ready(function(){
  url = window.location.pathname;
  urlRegExp = new RegExp(url.replace(/\/$/,'') + "$");
  if (url === '/'){
    $(".navbar-nav #dashboard").addClass("selected");
  }
  else{
    $(".navbar-nav li").each(function(){
      if(urlRegExp.test($(this).find('a').attr('href'))){
        $(this).addClass('selected');
      }
    });
  }
});