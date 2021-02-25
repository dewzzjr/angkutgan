const User = {
  Form: {},
  Set: function (data) {
    User.Form = data;
  },
  PassForm: {},
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
  Edit: function (callback, failedCallback) {
    let data = User.Form;
    $.ajax({
      type: 'PATCH',
      url: '/user/profile',
      contentType: 'application/json',
      data: JSON.stringify(data),
      success: function (data, status, xhr) {
        if (status === 'success') {
          User.Set(data.result);
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
  Validate: function (callback) {
    let ok = {
      message: [],
      valid: true
    };
    let data = User.Form;
    let check = function () {
      if (!data.fullname) {
        ok.valid = false;
        ok.message.push({
          name: "fullname",
          text: "nama tidak boleh kosong"
        });
      }
      if (callback) {
        callback(ok);
      }
      User.Set(data);
    }
    check();
    return ok
  },
  ValidatePassword: function (callback) {
    let ok = {
      message: [],
      valid: true
    };
    let data = User.PassForm;
    let check = function () {
      if (!data.old) {
        ok.valid = false;
        ok.message.push({
          name: "old",
          text: "password lama kosong"
        });
      }
      if (!data.new) {
        ok.valid = false;
        ok.message.push({
          name: "new",
          text: "password baru kosong"
        });
      }
      if (!data.repeat || data.repeat != data.new) {
        ok.valid = false;
        ok.message.push({
          name: "repeat",
          text: "password tidak sama"
        });
      }
      if (data.old == data.new) {
        ok.valid = false;
        ok.message.push({
          name: "new",
          text: "password tidak boleh sama password sebelumnya"
        });
      }
      if (callback) {
        callback(ok);
      }
      User.PassForm = data;
    }
    check();
    return ok
  },
  ChangePassword: function(callback, failedCallback) {
    let data = User.PassForm;
    $.ajax({
      type: 'POST',
      url: '/user/password',
      contentType: 'application/json',
      data: JSON.stringify(data),
      success: function (data, status, xhr) {
        if (status === 'success') {
          User.PassForm = {};
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
}
$(document).ready(function () {
  header((u) => {
    $('[name="username"]').val(u.username);
  });
  $('.datepicker').datepicker({
    format: 'dd/mm/yyyy'
  });
  User.GetProfile((data) => {
    $('[name="fullname"]').val(data.fullname);
    $('[name="nik"]').val(data.nik);
    $('[name="ktp"]').val(data.ktp);
    $('[name="birthdate"]').datepicker('update', data.birthdate);
    $('[name="religion"]').val(data.religion);
  });
  $('#formEdit').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let ubah = $(this).serializeObject();
    User.Set(ubah);
    User.Validate((ok) => {
      Form.Reset($(this));
      if (ok.valid) {
        Loading.Start($('#formEdit [type="submit"]'));
        User.Edit(() => {
          Form.Message('success', 'berhasil mengubah profil', $('#messageEdit'));
          Loading.End();
          Form.Reset($('#formEdit'));
        }, () => {
          Form.Message('danger', 'gagal mengubah profil', $('#messageEdit'));
          Loading.End();
        });
      } else {
        Form.Validate(ok.message, $('#formEdit'));
      }
    });
  });
  $('#formPass').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let pass = $(this).serializeObject();
    User.PassForm = pass;
    User.ValidatePassword((ok) => {
      Form.Reset($(this));
      if (ok.valid) {
        Loading.Start($('#formPass [type="submit"]'));
        User.ChangePassword(() => {
          Form.Message('success', 'berhasil mengubah password', $('#messagePass'));
          Loading.End();
          Form.Reset($('#formPass'));
        }, () => {
          Form.Message('danger', 'gagal mengubah password', $('#messagePass'));
          Loading.End();
        });
      } else {
        Form.Validate(ok.message, $('#formPass'));
      }
    });
  })
});