$(document).ready(function () {
  toastr.options = {
    "positionClass": "toast-bottom-right"
  };

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
    notification(message.type, message.content);
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
      var ackButton;
      async.series([
        // Get span color
        function(callback){
          getStyle(client.status, function(result){
            status = result;
            callback();
          });
        },
        function(callback){
          
          client.isSilenced(function(result){
            if(result){
              ackButton = "<span class='pull-right'><i class='fa fa-volume-off' onclick=\"deleteStash(event, '"+ client.name +"')\"></i></span>";
            }
            else {
              ackButton = "<span class='pull-right'><i class='fa fa-volume-up' onclick=\"postStash(event, '"+ client.name +"')\"></i></span>";
            }
            
          });
          callback();
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
              + "<span class='lead'>"+ client.name + ackButton +"  </span>"
              + "<span class='subtitle'><strong>"+ client.eventsCount() +"</strong></span>"
              + "<span class='small timestamp'><i class='fa fa-clock-o'></i> "+ client.last_check +"</span>"
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
        // Display clients
        clientsList.html(spans);
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
    //var ackButton = "";

    // Update title with client status
    client.isSilenced(function(result){
      if(result){
        var ackButton = "<i class='fa fa-volume-off'></i> <span class='pull-right' onclick=\"deleteStash(event, '"+ client.name +"')\"><button type='button' class='btn btn-danger btn-sm '>Un-silence client</button></span>";
      }
      else {
        var ackButton = "<i class='fa fa-volume-up'></i> <span class='pull-right' onclick=\"postStash(event, '"+ client.name +"')\"><button type='button' class='btn btn-sm btn-black'>Silence client</button></span>";
      }
      $("#client-details #name").html(ackButton + client.name);
    });

    // Update client details
    $("#client-details #address").html(client.address);
    $("#client-details #subscriptions").html(client.subscriptions.join(', '));
    $("#client-details #last-check").html(client.last_check);

    var parseHistory = function(data, nextCheck) {
      clientHistory = new History(data);
      var output = "";
      var eventDetails = "";
      var events = "";
      async.series([
        function(callback){
          clientHistory.getStyle(function(result){
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
          clientHistory.getEvent(clientEvents, client.name, function(err, result){
            if(err) console.log("Could not find active events for " + client.name);
            event = result;
            callback(err);
          });
        },
        function(callback){
          clientHistory.getCheck(function(err, result){
            if(err) console.log("Could not find the check " + clientHistory.check + " for client " + client.name);
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
          if($("td #"+client.name+"-"+clientHistory.check).hasClass('in')){
            detailsClass = "in";
          }
          else {
            detailsClass = "collapse";
          }
          callback();
        }
      ], function(err){
        if (!err){
          spans += "<tr data-toggle='collapse' data-target='#"+ client.name+"-"+clientHistory.check +"' class='accordion-toggle'>";

          // Status
          if (clientHistory.last_status == 0){
            if(clientHistory.last_execution == 0){
              spans += "<td><span class='label label-warning'>Inactive</span></td>";
            }
            else {
              spans += "<td><span class='label label-success'>Active</span></td>";
            }
          }
          else if(clientHistory.last_status == 1){
            spans += "<td><span class='label label-warning'>Warning</span></td>";
          }
          else if(clientHistory.last_status == 2){
            spans += "<td><span class='label label-danger'>Critical</span></td>";
          }
          else {
            spans += "<td><span class='label label-default'>Unknown</span></td>";
          }

          // Check name
          spans += "<td>"+ clientHistory.check +"</td>";

          // Output
          spans += (_.isObject(event)) ? "<td><span class='output'>"+ event.output +"</span></td>" : "<td></td>" ;
       
          // Last execution
          spans += "<td><i class='fa fa-clock-o'></i> "+ clientHistory.last_check +"</td>"
                  + "<td class='text-center'>";

          // Silence
          if(_.isObject(check)){
            //console.log('check is not null '+check);
            check.isSilenced("silence/"+client.name+"/"+check.name, function(result){
              if(result){
                spans += "<a href='#' class='btn btn-success btn-xs btn-hover' onclick=\"deleteStash(event, '"+ client.name +"', '"+ clientHistory.check +"')\"> <i class='fa fa-volume-off'></i></a>";
              }
              else {
                spans += "<a href='#' class='btn btn-warning btn-xs btn-hover' onclick=\"postStash(event, '"+ client.name +"', '"+ clientHistory.check +"')\"> <i class='fa fa-volume-up'></i></a>";
              }
            });
          }

          // Resolve
          if(_.isObject(event)){
            spans += "<a href='#' class='btn btn-danger btn-xs btn-hover' onclick=\"resolveEvent('"+ client.name +"', '"+ clientHistory.check +"')\"> <i class='fa fa-check'></i></a>";
          } else {
            spans += "<a href='#' class='btn btn-xs disabled'> <i class='fa fa-check'></i></a>";
          }

          // Events & Check details
          spans += "</td>"
                + "</tr>"
                + "<tr>"
                + "<td colspan='6' class='hiddenRow'><div id='"+ client.name+"-"+clientHistory.check +"' class='accordian-body "+ detailsClass +"'>";
          
          // Event details
          if(_.isObject(event)){
            spans += "<h5 class='title'><i class='fa fa-exclamation-triangle'></i> Event details</h5>"
              + "<dl class='dl-horizontal'>"
                + "<dt>Full output</dt>"
                + "<dd>"+event.output+"</dd>"
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
          spans += "<h5 class='title'><i class='fa fa-terminal'></i> Check details</h5>"
                + "<dl class='dl-horizontal'>"
                  + "<dt>Command</dt>"
                  + "<dd>"+check.command+"</dd>"
                  + "<dt>History</dt>"
                  + "<dd>"+clientHistory.history+"</dd>"
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