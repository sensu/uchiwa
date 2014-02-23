$(document).ready(function () {
  socket = io.connect();
  socket.emit('get_clients');

  // Clients
  socket.on('clients', function(data) {

    clients = JSON.parse(data.content);
    var spans = "";
    var clientsList = $("#clients-list");

    var parseClient = function(client, nextClient){

      var style = "block";
      var checks;
      var status;
      var subscriptions = "";

      async.series([
        // Get span color
        function(callback){
          getStyle(client.status, function(result){
            status = result;
            callback();
          });
        },
        function(callback){
          countEvents(client.events, function(result){
            checks = result;
            callback();
          });
        },
        // Is page filtered?
        function(callback){
          if(filter.clients && status == 'success'){
            style = "none";
          }
          callback();
        }
      ], function(err){
        if (err) return console.error("Error while fetching clients list: " + err);
        spans += "<a href='#' id='"+ client['name'] +"' class='list-group-item "+ status +"'  style='display: "+ style +";' data-toggle='modal' data-target='#client-details'><span class='name' style='min-width: 160px; display: inline-block;'><strong>"+ client['name'] +"</strong></span><span class=''>"+ checks +"</span><span class='text-muted' style='font-size: 12px;'></span><span class='badge'>"+ client['lastCheck'] +" ago</span><span class='pull-right'><i class='fa fa-clock-o'></i></span></a>";
        nextClient();
      });          
    };

    if(!$('#client-details').hasClass('in') && clientsList.length){
      // Parse each client to get the HTML span element
      async.each(clients, function(client, callback){
        parseClient(client, callback);
      },
      function(err){
        // Once we parsed each clients
        async.series([
          // Display message if no events found
          function(callback){
            var style;
            //var currentEvents = clientsList.find(".danger, .warning");
            if(events.length == 0) { // Do we have at least one alert?
              status = (filter.clients) ? "block" : "none";
              spans += "<span class='not-found' style='display: "+ status +";'><i class='fa fa-thumbs-o-up'></i> <h3>No alerts... for now!</h3></span>";
            }
            callback();
          },
          // Display last update message
          function(callback){
            var today = new Date();
            $('.last-update').html("<i class='fa fa-refresh fa-spin'></i> <small>" + today.toLocaleString() + "</small>");
            callback();
          }
        ], function(err){
          // Display HTML
          clientsList.html(spans);
        });
      });
    }
  });

  // Events
  socket.on('events', function(data) {
    events = JSON.parse(data.content);
    var spans = "";
    var eventsList = $("#events-list");

    var parseEvent = function(event, nextEvent){
      var status;
      var output;
      async.series([
        // Get status of the event
        function(callback){
          getStyle(event.status, function(result){
            status = result;
            callback();
          });
        },
        // Format the output
        function(callback){
          var maxLength = 40;
          output = event['output'];
          if(output.length > maxLength){
            output = output.substring(0,maxLength);
            output += "...";            
          }
          callback();
        }
      ], function(err){
        if (err) return console.error("Error while fetching events list: " + err);
        spans += "<a href='#' class='list-group-item "+ status +"'><span class='name' style='min-width: 160px; display: inline-block;'><strong>"+ event['client'] +"</strong></span><span class=''>"+ event['check'] +"</span><span class='text-muted' style='font-size: 12px;'> - "+ output +"</span><span class='badge'>"+ event['lastCheck'] +" ago</span><span class='pull-right'><i class='fa fa-clock-o'></i></span></a>";
        nextEvent();
      });
    };

    var displayEvents = function(){
      eventsList.html(spans);
    };

    if(!$('#event-details').hasClass('in') && eventsList.length){
      async.each(events, function(event, callback){
        parseEvent(event, callback);
      },
      function(err){
        // Once we parsed each events
        async.series([
          // Display message if no events found
          function(callback){
            var style;
            console.log(spans.length);
            //var currentEvents = clientsList.find(".danger, .warning");
            if(events.length == 0) { // Do we have at least one alert?
              status = (filter.clients) ? "block" : "none";
              spans += "<span class='not-found' style='display: "+ status +";'><i class='fa fa-thumbs-o-up'></i> <h3>No events... for now!</h3></span>";
            }
            callback();
          },
          // Display last update message
          function(callback){
            var today = new Date();
            $('.last-update').html("<i class='fa fa-refresh fa-spin'></i> <small>" + today.toLocaleString() + "</small>");
            callback();
          }
        ], function(err){
          // Display HTML
          eventsList.html(spans);
        });
      });
    }

  });

  // Client
  socket.on('client', function(data) {
  
    var checks = JSON.parse(data.content);
    var spans = "";
    var clientDetails =  $("#client-details");

    // Set name
    $("#client-details #name").html(currentClient);
 
    var parseCheck = function(check, nextCheck) {
      var output = "";
      var occurrences = "";
      async.series([
        function(callback){
          getStyle(check.last_status, function(result){
            status = result;
            callback();
          });
        },
        function(callback){
          if(check.last_status != 0 ) {
            findEvent(check.check, currentClient, function(result){
              output = result[0].output;
              occurrences = "- " + result[0].occurrences + " occurrence(s)";
            });
          }
          callback();
        },
        function(callback){
          var maxLength = 65;
          if(output.length > maxLength){
            output = output.substring(0,maxLength);
            output += "...";            
          }
          callback();
        },
        function(callback){
          if($("#checks #"+check['check']).hasClass('in')){
            detailsClass = "in";
          }
          else {
            detailsClass = "collapse";
          }
          callback();
        }
      ], function(err){
        if (err) return console.error("Error while fetching checks list: " + err);
        spans += "<a href='#' class='list-group-item "+ status +"' data-toggle='collapse' data-target='#"+ check.check + "'><span class='name' style='min-width: 180px; display: inline-block;'><strong>"+ check.check +"</strong></span><span class=''></span>"+ output +"<span class='text-muted' style='font-size: 12px;'> "+ occurrences +"</span><span class='badge'>"+ check.lastCheck +" ago</span><span class='pull-right'><i class='fa fa-clock-o'></i></span></a>";
        spans += "<span id='"+ check['check'] + "' class='"+ detailsClass + "'>"
          + "<div class='row'>"
            + "<div class='col-xs-6'>"
              + "<ul class='list-group'>"
                + "<li class='list-group-item'><strong>Full output</strong><span class='pull-right'><em>"+ check.output +"</em></span></li>"
                + "<li class='list-group-item'><strong>Last results</strong><span class='pull-right'><em>"+ check.history +"</em></span></li>"
              + "</ul>"
            + "</div>"
            + "<div class='col-xs-6'>"
            + "</div>"
          + "</div>"
          + "</span>";
        nextCheck();
      });
      
    }

    var displayChecks = function(){
      $("#client-details #checks").html(spans);
    };

    async.each(checks, function(check, callback){
      parseCheck(check, callback);
    },
    function(err){
      if (err) return console.error("Error while processing checks data: " + err);
      displayChecks();
    });

  });

  $("#clients-list").on('click', 'a', function(e) {
    getClient(this.id);
  });

});

var filter =  {
  clients: true,
  allClients: function(action){
    console.log(action);
    if (this.clients != action){
      this.clients = !this.clients;
      if(action) {
        $('.success').css("display", "none");
        $('.not-found').css("display", "block");
      }
      else {
        $('.success').css("display", "block");
        $('.not-found').css("display", "none");
      }
    }
  }
};

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

var findEvent = function(check, client, callback){
  var eventDetails = events.filter(function (e) { return e.client == client && e.check == check });
  callback(eventDetails);
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
  currentClient = id;
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