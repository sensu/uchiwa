var updateClients = function(clients) {

  var lines = [ "<div class='row'>" ];
  var list = $("div#clients-list");
  var template = [
    "<% if(clients.length == 0){ %>",
      "<div class='not-found'><i class='fa fa-exclamation-triangle'></i></div> <div class='not-found'>No clients found!</div>",
    "<% } %>",
    "<% var i = 1; %>",
    "<% _.each(clients, function(data) { %>",
      "<% var client = new Client(data); %>",
      "<% if(i % 4 === 0){ %>",
        "</div><div class='row'>",
      "<% } %>",
      "<% var isSilenced = client.isSilenced(); %>",
          "<div class='col-md-3 client'>",
            "<a href='#' id='<%= client.name %>' data-toggle='modal' data-target='#client-details'>",
              "<div class='well <%= getStyle(client.status) %>'>",
                "<span class='lead'><%= client.name %><span class='pull-right'><i class='fa fa-volume-<% isSilenced ? print('off') : print('up'); %>' onclick=\"<% isSilenced ? print('deleteStash') : print('postStash'); %>(event, '<%= client.name %>')\"></i></span></span>",
                "<span class='subtitle'><strong><%= client.eventsCount() %></strong></span>",
                "<span class='small timestamp'><i class='fa fa-clock-o'></i> <%= client.last_check %></span>",
              "</div>",
            "</a>",
          "</div>",
        "</div>",
      "<% i++; %>",
    "<% }); %>"
  ].join("");

  var lines = _.template(template, {clients: clients});
  if(!$('#client-details').hasClass('in') && list.length){
    list.html(lines);
  }

};
