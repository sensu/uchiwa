$(document).ready(function () {
  socket = io.connect();
  socket.emit('get_clients');

  socket.on('checks', function(data) {
    checks = JSON.parse(data.content);
  });

  socket.on('stashes', function(data) {
    stashes = JSON.parse(data.content);
  });

  socket.on('messenger', function(data) {
    var message = JSON.parse(data.content);
    displayMessage(message.type, message.page, message.content);
  });

  $("#clients-list").on('click', 'a', function(e) {
    getClient(this.id);
  });

  //
  // Clients
  //

  socket.on('clients', function(data) {

    clients = JSON.parse(data.content);
    var spans = "<div class='row'>";
    var clientsList = $("#clients-list");
    var i = 0;

    var parseClient = function(data, nextClient){
      client = new Client(data);
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
        }
      ], function(err){
        if (!err){
          if(i % 4 === 0){
            spans += "</div>";
            spans += "<div class='row'>";          
          }
          spans += ""
          + "<div class='col-md-3 client'>"
            + "<a href='#' id='"+ client.name +"' data-toggle='modal' data-target='#client-details'>"
            + "<div class='well "+ status +"'>"
              + "<span class='lead'>"+ client.name +"</span>"
              + "<span><strong>"+ client.eventsCount() +"</strong></span>"
              + "<span class='small'><i class='fa fa-clock-o'></i> "+ client.last_check +"</span>"
              + "</a>"
            + "</div>"
          + "</div>";
         i++;
        }        
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
  // Client details
  //

  socket.on('client', function(data) {
  
    client.history = JSON.parse(data.content);
    var spans = "";

    $("#client-details #name").html(client.name);
 
    var parseHistory = function(data, nextCheck) {
      history = new History(data);
      var output = "";
      var eventDetails = "";
      var events = "";
      async.series([
        function(callback){
          history.getStyle(function(result){
            status = result;
            callback();
          });
        },
        function(callback){
          client.getEvents(function(err, result){
            clientEvents = result;
            callback();
          });
        },
        function(callback){
          history.getEvent(clientEvents, client.name, function(err, result){
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
          if($("td #"+client.name+"-"+history.check).hasClass('in')){
            detailsClass = "in";
          }
          else {
            detailsClass = "collapse";
          }
          callback();
        }
      ], function(err){
        if (!err){
          spans += "<tr data-toggle='collapse' data-target='#"+ client.name+"-"+history.check +"' class='accordion-toggle'>";

          // Status
          if (history.last_status == 0){
            if(history.last_execution == 0){
              spans += "<td><span class='label label-warning'>Inactive</span></td>";
            }
            else {
              spans += "<td><span class='label label-success'>Active</span></td>";
            }
          }
          else if(history.last_status == 1){
            spans += "<td><span class='label label-warning'>Warning</span></td>";
          }
          else if(history.last_status == 2){
            spans += "<td><span class='label label-danger'>Critical</span></td>";
          }
          else {
            spans += "<td><span class='label label-default'>Unknown</span></td>";
          }

          // Check name
          spans += "<td>"+ history.check +"</td>";

          // Output
          spans += (_.isObject(event)) ? "<td>"+ event.output +"</td>" : "<td></td>" ;
       
          // Last execution
          spans += "<td><i class='fa fa-clock-o'></i> "+ history.last_check +"</td>"
                  + "<td class='text-center'>";

          // Silence
          if(_.isObject(check)){
            //console.log('check is not null '+check);
            check.isSilenced("silence/"+client.name+"/"+check.name, function(result){
              if(result){
                spans += "<a href='#' class='btn btn-success btn-xs btn-hover' onclick=\"deleteStash('"+ client.name +"', '"+ history.check +"')\"> <i class='fa fa-volume-off'></i></a>";
              }
              else {
                spans += "<a href='#' class='btn btn-warning btn-xs btn-hover' onclick=\"postStash('"+ client.name +"', '"+ history.check +"')\"> <i class='fa fa-volume-up'></i></a>";
              }
            });
          }

          // Resolve
          if(_.isObject(event)){
            spans += "<a href='#' class='btn btn-danger btn-xs btn-hover' onclick=\"resolveEvent('"+ client.name +"', '"+ history.check +"')\"> <i class='fa fa-check'></i></a>";
          } else {
            spans += "<a href='#' class='btn btn-xs disabled'> <i class='fa fa-check'></i></a>";
          }

          // Events & Check details
          spans += "</td>"
                + "</tr>"
                + "<tr>"
                + "<td colspan='6' class='hiddenRow'><div id='"+ client.name+"-"+history.check +"' class='accordian-body "+ detailsClass +"'>";
          
          // Event details
          if(_.isObject(event)){
            spans += "<h5>Event details</h5>"
              + "<dl class='dl-horizontal'>"
                + "<dt>Occurrences</dt>"
                + "<dd>"+event.occurrences+"</dd>"
                + "<dt>Flapping</dt>"
                + "<dd>"+event.flapping+"</dd>"
                + "<dt>Handlers</dt>"
                + "<dd>"+event.handlers+"</dd>"
                + "<dt>Issued</dt>"
                + "<dd>"+event.last_issued+"</dd>"
              + "</dl>";
          }
          
          // Check details
          spans += "<h5>Check details</h5>"
                + "<dl class='dl-horizontal'>"
                  + "<dt>Command</dt>"
                  + "<dd>"+check.command+"</dd>"
                  + "<dt>History</dt>"
                  + "<dd>"+history.history+"</dd>"
                  + "<dt>Subscribers</dt>"
                  + "<dd>"+check.subscribers+"</dd>"
                  + "<dt>Handlers</dt>"
                  + "<dd>"+check.handlers+"</dd>"
                  + "<dt>Interval</dt>"
                  + "<dd>"+check.interval+"</dd>"
                  + "<dt>Type</dt>"
                  + "<dd>"+check.type+"</dd>"
                  + "<dt>Handle</dt>"
                  + "<dd>"+check.handle+"</dd>"
                  + "<dt>Standalone</dt>"
                  + "<dd>"+check.standalone+"</dd>"
                  + "<dt>Subdue</dt>"
                  + "<dd>Begin: "+check.subdue.begin+" End: "+check.subdue.end+"</dd>";

          spans += "</dl>"; 
                + "</div></td>"
                + "</tr>";
        }
        
        nextCheck();
      });
    }

    async.each(client.history, function(check, callback){
      parseHistory(check, callback);
    },
    function(err){
      if (err) return console.error("Error while processing checks data: " + err);
      $("#client-details #historyList tbody").html(spans);
    });
  });
});