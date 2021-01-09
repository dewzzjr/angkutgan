const Pelanggan = {
  Form: {},
  Set: function (data) {
    this.Form = data;
  },
  Clear: function (callback) {
    this.Form = {};
    if (callback) {
      callback();
    }
  },
  Validate: function (isEdit = false, callback) {
    let ok = {
      message: [],
      valid: true
    };
    let data = this.Form;
    let set = this.Set;
    let check = function () {
      if (!data.code) {
        ok.valid = false;
        ok.message.push({
          name: "code",
          text: "kode tidak boleh kosong"
        });
      }
      if (!data.name) {
        ok.valid = false;
        ok.message.push({
          name: "name",
          text: "nama tidak boleh kosong"
        });
      }
      if (!data.type) {
        ok.valid = false;
        ok.message.push({
          name: "type",
          text: "pilih jenis pelanggan"
        });
      }
      if (!data.phone) {
        ok.valid = false;
        ok.message.push({
          name: "phone",
          text: "nomor telepon tidak boleh kosong"
        });
      }
      if (!data.address) {
        ok.valid = false;
        ok.message.push({
          name: "address",
          text: "alamat tidak boleh kosong"
        });
      }
      if (callback) {
        callback(ok);
      }
      set(data);
    }
    if (!isEdit) {
      $.ajax({
        type: 'GET',
        url: '/ajax',
        contentType: 'application/json',
        data: {
          action: 'validate_code_customer',
          new: data.code
        },
        success: function (data, status, xhr) {
          if (status === 'success' && !data.valid) {
            ok.valid = false;
            ok.message.push({
              name: "code",
              text: data.message
            });
          }
        },
        error: function (xhr, status, error) {
          console.log(status, error);
        }
      }).done(check);
    } else {
      check();
    }
    return ok
  },
  Create: function (callback, failedCallback) {
    let data = this.Form;
    let set = this.Set;
    $.ajax({
      type: 'POST',
      url: '/customer',
      contentType: 'application/json',
      data: JSON.stringify(data),
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
  Edit: function (callback, failedCallback) {
    let data = this.Form;
    let set = this.Set;
    $.ajax({
      type: 'PATCH',
      url: `/customer/${data.code}`,
      contentType: 'application/json',
      data: JSON.stringify(data),
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
  GetDetail: function (code, callback) {
    $.ajax({
      type: 'GET',
      url: `/customer/${code}`,
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success') {
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
  }
};

$(document).ready(function () {
  header();

  // TAMBAH
  $('#formAdd').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let tambah = $(this).serializeObject();
    console.log(tambah);
    Pelanggan.Set(tambah);
    Pelanggan.Validate(false, (ok) => {
      Form.Reset($('#formAdd'));
      if (ok.valid) {
        Loading.Start($('#formAdd [type="submit"]'));
        Pelanggan.Create(() => {
          Form.Message('success', 'berhasil menambah pelanggan', $('#messageAdd'));
          Loading.End();
          Form.Reset($('#formAdd'), () => {
            $('#formAdd input, #formAdd select').val('');
            $('#formAdd textarea').html('');
          });
        }, () => {
          Form.Message('danger', 'gagal menambah pelanggan', $('#messageAdd'));
          Loading.End();
        });
      } else {
        Form.Validate(ok.message, $('#formAdd'));
      }
    });
  });
});