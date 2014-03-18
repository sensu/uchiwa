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

var findClient = function(id, callback){
  var client = clients.filter(function (e) { return e.name == id });
  if(client.length > 0){
    callback(null, new Client(client[0]));
  }
  else {
    callback(true);
  }
}

/**
* getClient: Request to socket the client detail while modal window is shown
* @param id {string} Name of the client
*/
var getClient = function(id){
  var timer;
  clearInterval(timer);
  findClient(id, function(err, result){
    if (!err){
      client = result;
      socket.emit('get_client', {name: client.name});
      $('#client-details').on('hide.bs.modal', function () {
        $(this).off('hide.bs.modal');
        $("#client-details #checks").html("<span class='not-found'><i class='fa fa-spinner fa-spin'></i></span>");
        clearInterval(timer);
        client = null;
      })
      var timer = setInterval(function(){
        if($("#client-details").data('bs.modal').isShown){
          socket.emit('get_client', {name: client.name});
        }
      }, 10000);
    }
    else {
      console.log("Client '" + id + "' was not found.")
    }
  });
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

var postStash = function(id){
  var data = {path: id, content:{"reason": "uchiwa"}};
  socket.emit('create_stash', JSON.stringify(data));
};

var deleteStash = function(id){
  var data = {path: id, content:{}};
  socket.emit('delete_stash', JSON.stringify(data));
};

var displayMessage = function(type, page, message){
  if (type == "danger" || type == "warning" || type == "success" || type == "info"){
    type == "default";
  }

  var html = "<div class='alert alert-"+ type +" alert-dismissable>"
             + "<button type='button' class='close' data-dismiss='alert' aria-hidden='true'></button>"
             + message
             + "</div>";

  $("#" + page + " #message").html(html);
};