const User = {
  Form: {},
  GetProfile: function (callback) {
    $.ajax({
      type: 'GET',
      url: '/user/profile',
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success') {
          User.Form = data.result
          if (callback) {
            callback(data.result);
          }
        }
      },
      error: function (xhr, status, error) {
        console.log(status, error);
      },
    });
  },
}
$(document).ready(function () {
  header((u) => {
    $('[name="username"]').val(u.username);
  });
  User.GetProfile((data) => {

  })
});