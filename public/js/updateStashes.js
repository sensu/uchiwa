var updateStashes = function(stashes) {

  var lines = [ "<div class='row'>" ];
  var list = $("div#stashes-list");
  var template = [
    "<% if(stashes.length == 0){ %>",
      "<div class='not-found'><i class='fa fa-thumbs-o-up'></i></div> <div class='not-found'>No stashes found!</div>",
    "<% } %>",
    "<% var i = 1; %>",
    "<% _.each(stashes, function(data) { %>",
      "<% var stash = new Stash(data); %>",
      "<% if(i % 3 === 0){ %>",
        "</div><div class='row'>",
      "<% } %>",
        "<div class='col-md-4 client'>",
          "<div class='well default'>",
            "<span class='lead'><% (stash.client) ? print(stash.client) : print(stash.path); %><span class='pull-right'><i class='fa fa-times' onclick=\"deleteStash(event, '<%= stash.client %>', '<% (stash.check) ? print(stash.check) : print(''); %>')\"></i></span></span>",
            "<div class='row'>",
              "<div class='col-sm-5'>",
                "<span class='subtitle'><strong><% (stash.check) ? print(stash.check) : print('Client'); %></strong></span>",
                "<span class='small timestamp'><i class='fa fa-clock-o'></i> <%= stash.last_check %></span>",
              "</div>",
              "<div class='col-sm-7'>",
                "<% _.each(stash.content, function(value, key){ %>",
                  "<% if(key == 'timestamp') return; %>",
                  "<div class='small tag'><i class='fa fa-tag'></i> <strong><%= key %></strong> : <%= value %></div>",
                "<% }); %>",
              "</div>",
            "</div>",
          "</div>",
        "</div>",
      "</div>",
      "<% i++; %>",
    "<% }); %>"
  ].join("");

  var lines = _.template(template, {stashes: stashes});
  if(list.length){
    list.html(lines);
  }

};
