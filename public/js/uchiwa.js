$(document).ready(function(){

  toastr.options = {
    "positionClass": "toast-bottom-right"
  };

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

  $("#clients-list").on('click', 'a', function(e) {
    getClient(this.id);
  });
});

var getStyle = function(status){
  if (status == 2){
    return "danger";
  }
  else if (status == 1){
    return "warning";
  }
  else {
    return "success";
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

var postStash = function(e, client_name, check_name){
  var event = e || window.event;
  if (_.isUndefined(check_name)){
    check_name = "";
  }
  else {
    check_name = "/" + check_name;
  }
  var full_path = "silence/"+ client_name + check_name;
  var payload = {path: full_path, content:{"reason": "uchiwa"}};
  socket.emit('create_stash', JSON.stringify(payload));
  e.stopPropagation();
};

var deleteStash = function(e, client_name, check_name){
  var event = e || window.event;
  if (_.isUndefined(check_name)){
    check_name = "";
  }
  else {
    check_name = "/" + check_name;
  }
  var full_path = "silence/"+ client_name + check_name;
  var payload = {path: full_path, content:{}};
  socket.emit('delete_stash', JSON.stringify(payload));
  e.stopPropagation();
};

var resolveEvent = function(client_name, check_name){
  var payload = {client: client_name, check: check_name};
  socket.emit('resolve_event', JSON.stringify(payload));
};

var notification = function(type, message){
  toastr[type](message);
  fetch();
};

var fetch = function(){
  if(_.isString(client.name)){
    fetchAll();
    socket.emit('get_client', {name: client.name});
  }
};

var fetchAll = function(){
  async.series([
    function(callback){
      socket.emit('get_stashes', {});
      callback(null);
    },
    function(callback){
      socket.emit('get_checks', {});
      callback(null);
    },
    function(callback){
      socket.emit('get_events', {});
      callback(null);
    },
    function(callback){
      socket.emit('get_clients', {});
      callback(null);
    }
  ], function(err){
  });
}