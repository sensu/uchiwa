$(document).ready(function () {
  socket = io.connect();
  socket.emit('get_clients');

  //
  // Clients
  //

  socket.on('clients', function(data) {

    clients = JSON.parse(data.content);
    var spans = "<div class='row'>";
    var clientsList = $("#clients-list");
    var i = 0;

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
      ], function(err){
        if (!err){
          if(i % 4 === 0){
            spans += "</div>";
            spans += "<div class='row'>";          }
          spans += ""
          + "<div class='col-md-3 client'>"
            + "<a href='#' id='"+ client['name'] +"' data-toggle='modal' data-target='#client-details'>"
            + "<div class='well "+ status +"'>"
              + "<span class='lead'>"+ client['name'] +"</span>"
              + "<span><strong>"+ checks +"</strong></span>"
              + "<span class='small'><i class='fa fa-clock-o'></i> "+ client['last_check'] +"</span>"
              + "</a>"
            + "</div>"
          + "</div>";
         i++;
        }
        //spans += "<a href='#' id='"+ client['name'] +"' class='list-group-item "+ status +"' data-toggle='modal' data-target='#client-details'><span class='name' style='min-width: 160px; display: inline-block;'><strong>"+ client['name'] +"</strong></span><span class=''>"+ checks +"</span><span class='text-muted' style='font-size: 12px;'></span><span class='badge'>"+ client['last_check'] +" ago</span><span class='pull-right'><i class='fa fa-clock-o'></i></span></a>";
        
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
          }
        ], function(err){
          // Display HTML
          clientsList.html(spans);
        });
      });
    }
  });

  //
  // Events
  //

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
        spans += "<a href='#' class='list-group-item "+ status +"'><span class='name' style='min-width: 160px; display: inline-block;'><strong>"+ event['client'] +"</strong></span><span class=''>"+ event['check'] +"</span><span class='text-muted' style='font-size: 12px;'> - "+ output +"</span><span class='badge'>"+ event['last_check'] +" ago</span><span class='pull-right'><i class='fa fa-clock-o'></i></span></a>";
        nextEvent();
      });
    };

    if(!$('#event-details').hasClass('in') && eventsList.length){
      async.each(events, function(event, callback){
        parseEvent(event, callback);
      },
      function(err){
        async.series([
          function(callback){
            var style;
            if(events.length == 0) {
              status = (filter.clients) ? "block" : "none";
              spans += "<span class='not-found' style='display: "+ status +";'><i class='fa fa-thumbs-o-up'></i> <h3>No events... for now!</h3></span>";
            }
            callback();
          }
        ], function(err){
          // Display HTML
          eventsList.html(spans);
        });
      });
    }

  });

  //
  // Client
  //

  socket.on('client', function(data) {
  
    client.history = JSON.parse(data.content);
    var spans = "";
    var clientDetails = $("#client-details");

    $("#client-details #name").html(client.name);
 
    var parseCheck = function(data, nextCheck) {
      history = new History(data);
      var output = "";
      var eventDetails = "";
      async.series([
        function(callback){
          history.getStyle(function(result){
            status = result;
            callback();
          });
        },
        function(callback){
          history.getEvent(client.getEvents(), client.name, function(err, result){
            if(err) console.log("Could not find active events for " + client.name);
            event = result;
            callback(err);
          });
        },
        function(callback){
          history.getCheck(function(err, result){
            if(err) console.log("Could not find the check " + history.check + " for client " + client.name);
            check = result;
            callback(err);
          });
        },
        function(callback){
          if(typeof check.subscribers === undefined) {
            console.log("undefined");
            client.getSubscription(function(err, result){
              if(err) check.subscribers = "";
              check.subscribers = result;
              callback();
            });
          }
          else {
            callback();
          }
        },
        function(callback){
          if($("#checks #"+history['check']).hasClass('in')){
            detailsClass = "in";
          }
          else {
            detailsClass = "collapse";
          }
          callback();
        }
      ], function(err){
        if (!err){
          if(event) output = "<span class='output'>"+ event.output +"</span><span class='text-muted' style='font-size: 12px;'> - " + event.occurrences + " occurrences</span>";
          spans += "<a href='#' class='list-group-item "+ status +"' data-toggle='collapse' data-target='#"+ history.check + "'>"
            + "<span class='name' style='min-width: 180px; display: inline-block;'><strong>"+ history.check +"</strong></span>"
            + output
            + "<span class='badge'>"+ history.last_check +" ago</span><span class='pull-right'><i class='fa fa-clock-o'></i></span></a>";

          if(event) eventDetails = "<li class='list-group-item'><strong>Full output</strong><span class='pull-right'><em>"+ event.output +"</em></span></li>"
            + "<li class='list-group-item'><strong>Flapping</strong><span class='pull-right'><em>"+ event.flapping +"</em></span></li>"
            + "<li class='list-group-item'><strong>Event handlers</strong><span class='pull-right'><em>"+ event.handlers +"</em></span></li>";

          spans += "<span id='"+ history['check'] + "' class='"+ detailsClass + "'>"
            + "<div class='row'>"
              + "<div class='col-xs-12'>"
                + "<ul class='list-group'>"
                  + eventDetails
                  + "<li class='list-group-item'><strong>Last results</strong><span class='pull-right'><em>"+ history.history +"</em></span></li>"
                  + "<li class='list-group-item'><strong>Command</strong><span class='pull-right'><em>"+ check.command +"</em></span></li>"
                  + "<li class='list-group-item'><strong>Subscribers</strong><span class='pull-right'><em>"+ check.subscribers +"</em></span></li>"
                + "</ul>"
              + "</div>"
            + "</div>"
            + "</span>";
        }
        nextCheck();
      });
      
    }

    async.each(client.history, function(check, callback){
      parseCheck(check, callback);
    },
    function(err){
      if (err) return console.error("Error while processing checks data: " + err);
      $("#client-details #checks").html(spans);
      
      
    });

  });

  //
  // Checks
  //

  socket.on('checks', function(data) {
    checks = JSON.parse(data.content);
  });


  $("#clients-list").on('click', 'a', function(e) {
    getClient(this.id);
  });

});