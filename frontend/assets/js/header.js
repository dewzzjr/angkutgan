$.ajaxSetup({
  xhrFields: {
    withCredentials: true
  },
});

var Auth = {
  User: {},
  setUser: function (user) {
    this.User = user;
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
    let refresh = this.Refresh 
    if (this.User == {}) {
      window.location.replace('/login');
    } else {
      let set = this.setUser
      $.ajax({
        type: 'GET',
        url: 'user/info',
        contentType: 'application/json',
        success: function (data, status, xhr) {
          if (status === 'success') {
            set(data.result);
            if (callback) {
              callback(data.result);
            }
          }
        },
        error: function (xhr, status, error) {
          if (error == 'Unauthorized') {
            refresh();
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
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success') {
          set(data.result);
          if (callback) {
            callback(data.result);
          }
        }
      },
      error: function (xhr, status, error) {
        if (error == 'Unauthorized') {
          let tab = $(".tab-content .tab-pane.active.show").attr('id');
          let redirect = window.location.pathname + "?action=" + (tab ? tab : '');
          window.location.replace(`/login?redirect=${redirect}`);
        }
        console.log(status, error);
      },
    });
  },
};

function alert() {
  $('.alert').alert();
  $('.close').click(function(){
    $(this).parent().removeClass('show');
    $(this).parent().addClass('none');
 });
}

function header() {
  Auth.Get(function (u) {
    $('#headerUsername').html(u.username);
  });
  let path = window.location.pathname;
  let active = $(`[href="${path}"]`);
  active.addClass('active');
  active.attr('href', '#');
  
  $('#logout').click(function () {
    Auth.Logout(function () {
      console.log("OK")
      window.location.replace('/login')
    });
  });
  Menu.Init();
  alert();
}

function formatPrice(number, decPlaces = 2, decSep = ',', thouSep = '.') {
  decPlaces = isNaN(decPlaces = Math.abs(decPlaces)) ? 2 : decPlaces,
    decSep = typeof decSep === "undefined" ? "." : decSep;
  thouSep = typeof thouSep === "undefined" ? "," : thouSep;
  var sign = number < 0 ? "-" : "";
  var i = String(parseInt(number = Math.abs(Number(number) || 0).toFixed(decPlaces)));
  var j = (j = i.length) > 3 ? j % 3 : 0;

  return sign +
    (j ? i.substr(0, j) + thouSep : "") +
    i.substr(j).replace(/(\decSep{3})(?=\decSep)/g, "$1" + thouSep) +
    (decPlaces ? decSep + Math.abs(number - i).toFixed(decPlaces).slice(2) : "");
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
      $(input).parent('.form-group').append(`<div class="invalid-feedback">${e.text}</div>`);
    });
  },
};

const Menu = {
  Query: getUrlVars(),
  Init: function() {
    let action = this.Query['action']
    $(`#${action}`).tab('show');
    $(`[href="#${action}"]`).addClass('active');
  }
}