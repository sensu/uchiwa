var updateClient = function(client) {

  var titleTpl = [
    "<% var isSilenced = client.isSilenced(); %>",
    "<i class='fa fa-volume-<% isSilenced ? print('off') : print('up'); %>'></i>",
    "<span class='pull-right'onclick=\"<% isSilenced ? print('deleteStash') : print('postStash'); %>(event, '<%= client.name %>')\">",
      "<button type='button' class='btn btn-<% isSilenced ? print('danger') : print('black'); %> btn-sm '><% isSilenced ? print('Un-silence') : print('Silence'); %> client</button>",
    "</span>",
    " <%= client.name %>"
  ].join("");
  $("#client-details #name").html(_.template(titleTpl, {client: client}));

  $("#client-details #address").html(client.address);
  $("#client-details #subscriptions").html(client.subscriptions.join(', '));
  $("#client-details #last-check").html(client.last_check);

  var lines = [];

  async.each(client.history, function(data, nextCheck){
    checkHistory = new History(data);

    async.series([
      function(callback){
        client.getEvents(function(err, result){
          clientEvents = result;
          callback();
        });
      },
      function(callback){
        checkHistory.getEvent(clientEvents, client.name, function(err, result){
          if(err) console.log("Could not find active events for " + client.name);
          event = result;
          callback(err);
        });
      },
      function(callback){
        checkHistory.getCheck(function(err, result){
          if(err) console.log("Could not find the check " + checkHistory.check + " for client " + client.name);
          check = result;
          callback(err);
        });
      }
    ], function(err){
      if (!err){
        var checkTpl = [
          "<tr data-toggle='collapse' data-target='#<%= client.name %>-<%= checkHistory.check %>' class='accordion-toggle'>",
            "<% if(!checkHistory.last_status){ %>",
              "<td><span class='label label-<% checkHistory.last_execution ? print('success') : print('warning'); %>'><% checkHistory.last_execution ? print('Active') : print('Inactive'); %></span></td>",
            "<% } else { %>",
              "<td><span class='label label-<% (checkHistory.last_status == 1) ? print('warning') : print('danger'); %>'><% (checkHistory.last_status == 1) ? print('Warning') : print('Critical'); %></span></td>",
            "<% } %>",
            "<td><%= checkHistory.check %></td>",
            "<% _.isObject(event) ? print('<td><span class=\"output\">' + event.output + '</span></td>') : print('<td></td>'); %>",
            "<td><i class='fa fa-clock-o'></i> <%= checkHistory.last_check %></td>",
            "<td class='text-center'>",
              "<% var isSilenced = check.isSilenced('silence/'+ client.name +'/'+ checkHistory.check); %>",
              "<a href='#' class='btn btn-xs btn-hover btn-<% isSilenced ? print('success') : print('warning'); %>' onclick=\"<% isSilenced ? print('deleteStash') : print('postStash'); %>(event, '<%= client.name %>', '<%= checkHistory.check %>')\"> ",
              "<i class='fa fa-volume-<% isSilenced ? print('off') : print('up'); %>'></i></a>",
              "<% if(_.isObject(event)){ %>",
                "<a href='#' class='btn btn-danger btn-xs btn-hover' onclick=\"resolveEvent('<%= client.name %>', '<%= checkHistory.check %>')\"> <i class='fa fa-check'></i></a>",
              "<% } else { %>",
                "<a href='#' class='btn btn-xs disabled'> <i class='fa fa-check'></i></a>",
              "<% } %>",
            "</td>",
          "</tr>",
          "<tr>",
           "<% var isCollapsed = $('td #'+client.name+'-'+checkHistory.check).hasClass('in') %>",
            "<td colspan='6' class='hiddenRow'>",
              "<div id='<%= client.name %>-<%= checkHistory.check %>' class='accordian-body <% (isCollapsed) ? print('in') : print('collapse'); %>'>",
                "<% if(_.isObject(event)){ %>",
                  "<h5 class='title'><i class='fa fa-exclamation-triangle'></i> Event details</h5>",
                  "<dl class='dl-horizontal'>",
                    "<dt>Full output</dt>",
                    "<dd><%= event.output %></dd>",
                    "<dt>Occurrences</dt>",
                    "<dd><%= event.occurrences %></dd>",
                    "<dt>Flapping</dt>",
                    "<dd><%= event.flapping %></dd>",
                    "<dt>Handlers</dt>",
                    "<dd><%= event.handlers %></dd>",
                    "<dt>Issued</dt>",
                    "<dd><%= event.last_issued %></dd>",
                  "</dl>",
                "<% } %>",
                "<h5 class='title'><i class='fa fa-terminal'></i> Check details</h5>",
                "<dl class='dl-horizontal'>",
                    "<dt>Command</dt>",
                    "<dd><%= check.command %></dd>",
                    "<dt>History</dt>",
                    "<dd><%= checkHistory.history %></dd>",
                    "<dt>Subscribers</dt>",
                    "<dd><%= check.subscribers %></dd>",
                    "<dt>Handlers</dt>",
                    "<dd><%= check.handlers %></dd>",
                    "<dt>Interval</dt>",
                    "<dd><%= check.interval %></dd>",
                    "<dt>Type</dt>",
                    "<dd><%= check.type %></dd>",
                    "<dt>Handle</dt>",
                    "<dd><%= check.handle %></dd>",
                    "<dt>Standalone</dt>",
                    "<dd><%= check.standalone %></dd>",
                    "<dt>Subdue</dt>",
                    "<dd>Begin: <%= check.subdue.begin %> End: <%= check.subdue.end %></dd>",
                "</dl>",
              "</div>",
            "</td>",
          "</tr>",
        ].join("");
        lines.push(_.template(checkTpl, {checkHistory: checkHistory, client: client, event: event, check: check}));
      }
      nextCheck();
    });
  },
  function(err){
    if (err) return console.error("Error while processing checks data: " + err);
    $("#client-details #historyList tbody").html(lines);
  });

};