var updateChecks = function(checks) {

  var lines = [ "<div class='row'>" ];
  var list = $("div#checks-list");
  var template = [
    "<% if(checks.length == 0){ %>",
      "<div class='not-found'><i class='fa fa-exclamation-triangle'></i></div> <div class='not-found'>No checks found!</div>",
    "<% } %>",
    "<% var i = 1; %>",
    "<% _.each(checks, function(data) { %>",
      "<% var check = new Check(data); %>",
      "<% if(i % 2 === 0){ %>",
        "</div><div class='row'>",
      "<% } %>",
        "<div class='col-md-6 client'>",
          "<div class='well default'>",
            "<span class='lead'><%= check.name %><span class='pull-right small tag'><i class='fa fa-tag'></i> <%= check.subscribers %></span></span>",
            "<span class='subtitle'><i class='fa fa-terminal'></i> <strong><%= check.command %></strong></span>",
            "<span class='small timestamp'><i class='fa fa-clock-o'></i> <%= check.interval %> seconds</span>",
            "<span class='small'> <% (!check.standalone) ? print('Standalone') : print('Not standalone'); %></span>",
          "</div>",
        "</div>",
      "</div>",
      "<% i++; %>",
    "<% }); %>"
  ].join("");

  var lines = _.template(template, {checks: checks});
  if(list.length){
    list.html(lines);
  }

};
