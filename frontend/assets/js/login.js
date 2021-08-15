$(document).ready(function () {
  $('.alert').alert();
  $('.close').click(function () {
    $(this).parent().removeClass('show');
    $(this).parent().addClass('none');
  });
  $('#login').submit(function (e) {
    e.preventDefault();
    $('#login [type="submit"]').append('<span class="spinner-border spinner-border-sm" role="status"></span>');
    $('#login [type="submit"]').attr('disabled', 'disabled');
    let password = $('#password').val();
    let username = $('#username').val();
    Auth.Login(username, password, function () {
      let url = getUrlVars()['redirect'];
      url = decodeURIComponent(url);
      window.location.replace(url ? url : '/');
    }, function (e) {
      if (e == 'Unauthorized') {
        e = 'Username or password is wrong';
      }
      $('#message').html(e);
      $('.alert').removeClass('none');
      $('.alert').addClass('show');
      $('[role="status"]').remove();
      $('#login [type="submit"]').removeAttr('disabled');
    });
  });
});