$(document).ready(function(){
  socket = io.connect();

  /**
   * Notifications
   */
  toastr.options = {
    "positionClass": "toast-bottom-right"
  };
  var notification = function(type, message){
    toastr[type](message);
  };
  socket.on('messenger', function(data) {
    if(_.isUndefined(data.content)) return;
    var message = JSON.parse(data.content);
    notification(message.type, message.content);
  });

  /**
   * Navbar
   */
  $('.navbar-nav [data-toggle="tooltip"]').tooltip();
  $('.navbar-twitch-toggle').on('click', function(event) {
    event.preventDefault();
    $('.navbar-twitch').toggleClass('open');
  });

  /**
  * Graphics
  */
  dashboardGraph = Morris.Line({
      element: 'dashboard-graph',
      data: [],
      xkey: 'y',
      ykeys: ['e', 's'],
      labels: ['Events', 'Stashes'],
      lineColors: ['#2CA7E5', '#F9CD65'],
      hideHover: 'auto',
      pointSize: 0,
      fillOpacity: 1,
      gridTextColor: '#fff',
      gridTextFamily: "'Lato', sans-serif",
      gridTextWeight: 700,
      grid: false,
      lineWidth: 4,
      axes: true,
      behaveLikeLine: true
  });

  socket.on('stats', function(data) {
    if(_.isUndefined(data.content) || $('#dashboard-graph').length == 0) return;
    var stats = JSON.parse(data.content);
    dashboardGraph.setData(stats);
  });
});
