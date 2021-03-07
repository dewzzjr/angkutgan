
$.ajaxSetup({
  xhrFields: {
    withCredentials: true
  }
});

var Auth = {
  User: {},
  setUser: function (user) {
    this.User = user;
  },
  countdown: 0,
  unauthorized: function () {
    let tab = $(".tab-content .tab-pane.active.show").attr('id');
    let redirect = window.location.pathname + (tab ? "?action=" + tab : '');
    redirect = encodeURIComponent(redirect);
    window.location.replace(`/login?redirect=${redirect}`);
  },
  Login: function (username, password, callback, failedCallback) {
    let set = this.setUser
    $.ajax({
      type: 'POST',
      url: '/user/login',
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
          Auth.CheckExpire();
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
      url: '/user/logout',
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
      $.ajax({
        type: 'GET',
        url: '/user/info',
        contentType: 'application/json',
        beforeSend: function() {
          Auth.CheckExpire();
        },
        success: function (data, status, xhr) {
          if (status === 'success') {
            Auth.setUser(data.result);
            if (callback) {
              callback(data.result);
            }
            Auth.CheckExpire();
          }
        },
        error: function (xhr, status, error) {
          Auth.unauthorized();
          console.log(status, error);
        },
      });
    }
  },
  Refresh: function (callback, failedCallback) {
    $.ajax({
      type: 'POST',
      url: '/user/session',
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success') {
          Auth.setUser(data.result);
          if (callback) {
            callback(data.result);
          }
          Auth.CheckExpire();
        }
      },
      error: function (xhr, status, error) {
        if (failedCallback) {
          failedCallback()
        }
        if (xhr.status == 401) {
          Auth.unauthorized();
        }
        console.log(status, error);
      },
    });
  },
  CheckExpire: function (callback, failedCallback) {
    clearTimeout(Auth.countdown);
    Auth.countdown = setTimeout(() => {
      Auth.Refresh(callback, failedCallback)
    }, (Auth.User.exp - 30) * 1000 - $.now());
  }
};

const Daftar = {
  Page: 1,
  Rows: 10,
  Data: [],
  Keyword: '',
  URL: '',
  Init: function(url) {
    Daftar.URL = url;
  },
  SetData: function (data) {
    Daftar.Data = data;
  },
  CacheFunc: {},
  SetFunc: function (name, callback, failedCallback) {
    this.CacheFunc[name] = {
      callback: callback,
      failedCallback: failedCallback,
    };
  },
  GetData: function (callback, failedCallback) {
    let page = this.Page;
    let rows = this.Rows;
    let keyword = this.Keyword;
    let set = this.SetData;
    this.SetFunc('GetData', callback, failedCallback);
    $.ajax({
      type: 'GET',
      url: this.URL,
      contentType: 'application/json',
      data: {
        page: page,
        row: rows,
        keyword: keyword,
      },
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
  Apply: function(callback) {
    let func = this.CacheFunc['GetData']
    if (func) {
      this.GetData(function(data) {
        if (func.callback) {
          func.callback(data);
        }
        if (callback) {
          callback();
        }
      }, func.failedCallback);
    }
  },
  Reload: function() {
    Daftar.Page = 1;
    Daftar.Rows = 10;
    Daftar.Keyword = '';
    let func = this.CacheFunc['GetData']
    if (func) {
      this.GetData(func.callback, func.failedCallback);
    }
  }
};

function alert() {
  $('.alert').alert();
  $('.close').click(function(){
    $(this).parent().removeClass('show');
    $(this).parent().addClass('none');
 });
}

function header(callback) {
  Auth.Get(function (u) {
    $('#headerUsername').html(u.username);
    if (callback) {
      callback(u);
    }
  });
  let path = window.location.pathname;
  let active = $(`[href="${path}"]`);
  active.addClass('active');
  active.attr('href', '#');
  
  $('#logout').click(function () {
    Auth.Logout(function () {
      Auth.unauthorized();
    });
  });
  $('[name="code"]').on('keypress', function (e) {
      var regex = new RegExp("^[a-zA-Z0-9]+$");
      var str = String.fromCharCode(!e.charCode ? e.which : e.charCode);
      if (regex.test(str)) {
          return true;
      }
  
      e.preventDefault();
      return false;
  })
  Menu.Init();
  alert();
}

function formatPrice(number) {
  var formatter = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    currencySign: 'accounting',
  
    // These options are needed to round to whole numbers if that's what you want.
    minimumFractionDigits: 2, // (this suffices for whole numbers, but will print 2500.10 as $2,500.1)
    maximumFractionDigits: 2, // (causes 2500.99 to be printed as $2,501)
  });
  return formatter.format(number);
}

function isInt(value) {
  var x = parseFloat(value);
  return !isNaN(value) && (x | 0) === x;
}

function getUrlVars()
{
    var vars = [], hash;
    var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
    for(var i = 0; i < hashes.length; i++)
    {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
    }
    return vars;
}
const Loading = {
  Button: {},
  Spinner: {},
  Start: function(button) {
    this.Button = button;
    $(button).append('<span class="spinner-border spinner-border-sm" role="status"></span>');
    this.Spinner = $(this.Button).find('[role="status"]');
    $(button).attr('disabled', 'disabled');
  },
  End: function() {
    $(this.Spinner).remove();
    $(this.Button).removeAttr('disabled');
  },
};

const Form = {
  Component: {},
  Init: function(form, message) {
    if (!message) {
      message = $(form).find('.text-validation');
    }
    this.Component = {
      message: $(message),
      form: $(form),
    };
  },
  Message: function(tag, message, component) {
    if (!component) {
      component = $(this.Component.message);
    }
    $(component).html(message);
    $(component).parent().removeAttr('class');
    $(component).parent().addClass(`alert alert-${tag} fade show`);
    alert();
  },
  Reset: function (form, callback) {
    if (!form) {
      form = $(this.Component.form);
    }

    $(form).removeClass('was-validated');
    $(form).find('.invalid-feedback').remove();
    $(form).find('.is-invalid').removeClass('is-invalid');

    if (callback) {
      callback();
    }
  },
  Validate: function (message, form) {
    if (!form) {
      form = $(this.Component.form);
    }
    $(form).addClass('was-validated');
    $(form).find('.invalid-feedback').remove();
    console.log(message);

    message.forEach(e => {
      let input = $(form).find(`[name="${e.name}"]`);
      $(input).addClass('is-invalid');
      $(input).parents().find('.input-group').append(`<div class="invalid-feedback">${e.text}</div>`);
    });
  },
};

const Menu = {
  Query: getUrlVars(),
  Init: function() {
    let action = this.Query['action'];
    if (action) {
      $(`#${action}`).tab('show');
      $(`[href="#${action}"]`).addClass('active');
    }
  }
}