var Auth = {
  User: {},
  setUser: function (user) {
    this.User = user;
  },
  Login: function (username, password, callback, failedCallback) {
    let set = this.setUser
    $.ajax({
      type: 'POST',
      url: 'user/login',
      contentType: 'application/json',
      data: JSON.stringify({
        username: username,
        password: password,
      }),
      success: function (data, status, xhr) {
        if (status === 'success') {
          set(data.result);
          if (callback) {
            callback(data.result);
          }
        }
      },
      error: function (xhr, status, error) {
        console.log(status, error);
        if (failedCallback) {
          failedCallback(error);
        }
      },
    });
  },
  Logout: function (callback, failedCallback) {
    let set = this.setUser
    $.ajax({
      type: 'POST',
      url: 'user/logout',
      contentType: 'application/json',
      success: function (data, status, xhr) {
        set({});
        if (callback) {
          callback(data.result);
        }
      },
      error: function (xhr, status, error) {
        console.log(status, error);
        if (failedCallback) {
          failedCallback(error);
        }
      },
    });
  },
  Get: function (callback) {
    if (this.User == {}) {
      window.location.replace('/login');
    } else {
      let set = this.setUser
      $.ajax({
        type: 'GET',
        url: 'user/info',
        xhrFields: {
          withCredentials: true
        },
        contentType: 'application/json',
        success: function (data, status, xhr) {
          if (status === 'success') {
            set(data.result);
            if (callback) {
              // console.log(set, data.result);
              callback(data.result);
            }
          }
        },
        error: function (xhr, status, error) {
          if (error == 'Unauthorized') {
            window.location.replace('/login');
          }
          console.log(status, error);
        },
      });
    }
  },
  Refresh: function (callback) {
    let set = this.setUser
    $.ajax({
      type: 'GET',
      url: 'user/info',
      xhrFields: {
        withCredentials: true
      },
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success') {
          set(data.result);
          if (callback) {
            // console.log(set, data.result);
            callback(data.result);
          }
        }
      },
      error: function (xhr, status, error) {
        if (error == 'Unauthorized') {
          window.location.replace('/login');
        }
        console.log(status, error);
      },
    });
  },
};

function header() {
  Auth.Get(function (u) {
    $('#headerUsername').html(u.username);
  });
  let path = window.location.pathname;
  let active = $(`[href="${path}"]`);
  active.addClass('active');
  active.attr('href', '#');
  $('#logout').click(function () {
    Auth.Logout(function() {
      console.log("OK")
      window.location.replace('/login')
    });
  });
}