var Auth = {
  User: {},
  setUser: function (user) {
    this.User = user;
  },
  Login: function (username, password, callback) {
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
            // console.log(set, data.result);
            callback(data.result);
          }
        }
      },
      error: function (xhr, status, error) {
        if (error == 'Unauthorized') {
          // TODO
        }
        console.log(status, error);
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