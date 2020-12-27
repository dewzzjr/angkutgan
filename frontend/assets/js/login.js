$(document).ready(function () {
  $('.alert').alert();
  $('.close').click(function(){
    $(this).parent().removeClass('show');
 });
  $('#login').submit(function (e) {
    e.preventDefault();
    let password = $('#password').val();
    let username = $('#username').val();
    Auth.Login(username, password, function() {
      window.location.replace('/');
    }, function (e) {
      if (e == 'Unauthorized') {
        e = 'Username or password is wrong';
      }
      $('#message').html(e);
      $('.alert').addClass('show');
    });
  });
});