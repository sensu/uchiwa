var updateEvents = function(events) {

  var lines = [ "<div class='row'>" ];
  var list = $("div#events-list");
  var template = [
    "<% if(events.length == 0){ %>",
      "<div class='not-found'><i class='fa fa-thumbs-o-up'></i></div> <div class='not-found'>No events found!</div>",
    "<% } %>",
    "<% var i = 1; %>",
    "<% _.each(events, function(data) { %>",
      "<% var event = new Event(data); %>",
      "<% if(i % 4 === 0){ %>",
        "</div><div class='row'>",
      "<% } %>",
      "<%  %>",
          "<div class='col-md-3 client'>",
            "<a href='#' id='<%= event.client %>' data-toggle='modal' data-target='#client-details'>",
              "<div class='well <% (event.status == 1) ? print('warning') : print('danger'); %>'>",
                "<span class='lead'><%= event.client %><span class='pull-right'><i class='fa fa-volume-up'></i></span></span>",
                "<span class='subtitle'><strong><%= event.check %></strong></span>",
                "<span class='small timestamp'><i class='fa fa-clock-o'></i> <%= event.last_issued %></span>",
              "</div>",
            "</a>",
          "</div>",
        "</div>",
      "<% i++; %>",
    "<% }); %>"
  ].join("");

  var lines = _.template(template, {events: events});
  if(list.length){
    list.html(lines);
  }

};
