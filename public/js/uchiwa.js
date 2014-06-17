$(document).ready(function(){

  /**
   * Navbar
   */
  $('.navbar-nav [data-toggle="tooltip"]').tooltip();
  $('.navbar-twitch-toggle').on('click', function(event) {
    event.preventDefault();
    $('.navbar-twitch').toggleClass('open');
  });

});
