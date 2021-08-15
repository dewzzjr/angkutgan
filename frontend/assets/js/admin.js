const User = {
  Form: {
    username: '', //validate ajax di tambah user
    fullname: '',
    address: '',
    phone: '',
    nik: '',
    ktp: '',
    religion: '',
    birthdate: ''
  },
  Set: function (data) {
    User.Form = data;
  },
  Clear: function (callback) {
    User.Form = {};
    if (callback) {
      callback();
    }
  },
  Create: function (callback, failedCallback) {
    let data = this.Form;
    let set = this.Set;
    let retries = false;
    $.ajax({
      type: 'POST',
      url: '/user/create',
      contentType: 'application/json',
      data: JSON.stringify(data),
      success: function (data, status, xhr) {
        if (status === 'success') {
          set(data.result);
          if (callback) {
            callback(data.result);
            Daftar.Reload();
          }
        }
      },
      error: function (xhr, status, error) {
        console.log(status, error);
        if (failedCallback) {
          failedCallback(error);
        }
        if (xhr.status == 401 && !retries) {
          retries = true;
          Auth.Refresh(function() {
            $.ajax(this);
          });
        }
      },
    });
  }
};

$(document).ready(function () {
  // Example starter JavaScript for disabling form submissions if there are invalid fields
  
  
  header();
  var user = Auth.User;
  console.log(user);

  // Convert date format
  function dateFormat(date){
    var dateAr = date.split('-');
    var newDate = dateAr[2] + '/' + dateAr[1] + '/' + dateAr[0];
    return newDate;
  }

  // Reset form after submit
  function resetForm() {
      $('#codeAdd').val('');
      $('#nameAdd').attr('disabled', true).val('');
      $('#phoneAdd').attr('disabled', true).val('');
      $('#addressAdd').attr('disabled', true).val('');
      $('#nipAdd').attr('disabled', true).val('');
      $('#ktpAdd').attr('disabled', true).val('');
      $('#ttlAdd').attr('disabled', true).val('');
      $('#religionAdd').attr('disabled', true).val('');
  }

  function invalid($param){
    if (!$param.val()) {
      $param.addClass('is-invalid');
      setTimeout (function(){
        $param.removeClass('is-invalid');
      }, 3000);
    }
  }

  // Tambah user
  $('#codeAdd').on('change paste keyup', function() {
    $.ajax({
      type: 'GET',
      url: `/ajax?action=username&new=${$(this).val()}`,
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success' && data.valid) {
          $('#nameAdd').attr('disabled', false);
          $('#phoneAdd').attr('disabled', false);
          $('#addressAdd').attr('disabled', false).val('');
          $('#nipAdd').attr('disabled', false);
          $('#ktpAdd').attr('disabled', false);
          $('#ttlAdd').attr('disabled', false);
          $('#religionAdd').attr('disabled', false);
        }
        console.log(data, status);
      },
      error: function (xhr, status, error) {
        console.log(status, error);
        $.alert({
          title: 'Gagal',
          content: `${$(this).val()}: ${data.message}`,
        });
        invalid($(this));
      }
    });
  });

  $('.submit').on('click', function (e) {
    if (!$('#codeAdd').val() || !$('#nameAdd').val() || !$('#religionAdd').val() || !$('#ttlAdd').val()) {
      $('#alert').addClass("alert-danger").show();
      $('#alert #messageAdd').text("Lengkapi form sebelum submit");
      setTimeout (function(){
        $('#alert').removeClass("alert-danger").hide('fade');
      }, 3000);
      if (!$('#codeAdd').val()) {
        invalid($('#codeAdd'));
      }
      if (!$('#nameAdd').val()) {
        invalid($('#nameAdd'));
      }
      if (!$('#religionAdd').val()) {
        invalid($('#religionAdd'));
      }
      if (!$('#ttlAdd').val()) {
        invalid($('#ttlAdd'));
      }
    } else {
      User.Form ["username"] = $('#codeAdd').val();
      User.Form ["fullname"] = $('#nameAdd').val();
      User.Form ["address"] = $('#addressAdd').val();
      User.Form ["phone"] = $('#phoneAdd').val();
      User.Form ["nik"] = $('#nipAdd').val();
      User.Form ["ktp"] = $('#ktpAdd').val();
      User.Form ["religion"] = $('#religionAdd').val();
      User.Form ["birthdate"] = dateFormat($('#ttlAdd').val());

      User.Create(() => {
        $('#success').addClass("alert-success").show();
        $('#success #messageAdd').text("User berhasil ditambahkan");
        setTimeout (function(){
          $('#success').removeClass("alert-success").hide('fade');
        }, 3000);
        Loading.End();

        resetForm();
        $('body').scrollTop(0);

      }, () => {
        $('#alert').addClass("alert-danger").show();
        $('#alert #messageAdd').text("Username sudah digunakan, coba lagi");
        setTimeout (function(){
          $('#alert').removeClass("alert-danger").hide('fade');
        }, 3000);
        Loading.End();
      });

      console.log(User.Form);
    }

    e.preventDefault();
  });

  //Daftar user
  


});