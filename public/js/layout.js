// Add 'selected' class to active menu element
$(document).ready(function(){
  url = window.location.pathname;
  urlRegExp = new RegExp(url.replace(/\/$/,'') + "$");
  if (url === '/'){
    $(".navbar-nav #dashboard").addClass("selected");
  }
  else{
    $(".navbar-nav li").each(function(){
      if(urlRegExp.test($(this).find('a').attr('href'))){
        $(this).addClass('selected');
      }
    });
  }
});