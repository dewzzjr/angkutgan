const Pelanggan = {
  Form: {
    code: '',
    name: '',
    type: 0,
    phone: '',
    address: '',
    project: []
  },
  Set: function (data) {
    Pelanggan.Form = data;
  },
  Clear: function (callback) {
    Pelanggan.Form = {};
    if (callback) {
      callback();
    }
  },
  AppendProject: function (project) {
    Pelanggan.Form.project.push(project);
  },
  EditProject: function (index, project) {
    Pelanggan.Form.project[index] = project;
  },
  LenProject: function () {
    return this.Form.project.length;
  },
  RemoveProject: function (index, callback) {
    let pop = Pelanggan.Form.project.splice(index, 1);
    if (callback) {
      callback(pop);
    }
  },
  SetProject: function (form, func, before, after) {
    if (before) {
      before();
    }
    Pelanggan.Form.project.forEach(function (p, i) {
      console.log(i, p);
      func($(form), i, p);
    });
    if (after) {
      after();
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
      } else {
        data.code = data.code.toUpperCase();
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
      } else {
        data.type = parseInt(data.type);
      }

      if (data.type == 2) {
        if (!data.group_name) {
          ok.valid = false;
          ok.message.push({
            name: "group_name",
            text: "nama perusahaan tidak boleh kosong"
          });
        }
        if (!data.role) {
          ok.valid = false;
          ok.message.push({
            name: "role",
            text: "jabatan tidak boleh kosong"
          });
        }
        data.project.forEach((p, i) => {
          if (!p.name) {
            ok.valid = false;
            ok.message.push({
              name: `project[${i}][name]`,
              text: "nama proyek tidak boleh kosong"
            });
          }
          if (!p.location) {
            ok.valid = false;
            ok.message.push({
              name: `project[${i}][location]`,
              text: "lokasi proyek tidak boleh kosong"
            });
          }
        })
      } else {
        data.group_name = "";
        data.role = "";
        data.project = [];
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
            Daftar.Reload();
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
            Daftar.Reload();
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

  Daftar.Init('/customers');
  Daftar.GetData(function (data) {
    $('#tablePelanggan tbody').html('');
    data.forEach(e => {
      let empty = {
        name: '',
        value: ''
      };
      var nik = (e.nik) ? {
        name: `<p>NIK: </p>`,
        value: `<p>${e.nik}</p>`
      } : empty;
      var role = (e.role) ? {
        name: `<p>Sebagai: </p>`,
        value: `<p>${e.role}</p>`
      } : empty;
      var groupName = (e.group_name) ? {
        name: `<p>Perusahaan: </p>`,
        value: `<p>${e.group_name}</p>`
      } : empty;
      var name = (e.group_name) ? e.group_name : e.name;
      var row = `
      <tr>
        <th scope="row">${e.code}</th>
        <td>${name}</td>
        <td>${e.address}</td>
        <td>
            <div class="btn-group">
                <button type="button" class="btn btn-success" data-toggle="collapse" data-target="#${e.code}"
                    data-parent="#tablePelanggan" class="collapsed">Detail</button>
                <button type="button" class="btn btn-warning">Ubah</button>
                <button type="button" class="btn btn-danger">Hapus</button>
            </div>
        </td>
      </tr>
      <tr class="collapse" id="${e.code}">
        <th class="text-right" scope="row" colspan="2">
          <p>Jenis: </p>
          ${groupName.name}
          ${nik.name}
          <p>Nama: </p>
          ${role.name}
          <p>No. Telp: </p>
        </th>
        <td colspan="3">
          <p>${e.type_desc}</p>
          ${groupName.value}
          ${nik.value}
          <p>${e.name}</p>
          ${role.value}
          <p>${e.phone}</p>

          <button type="button" class="btn btn-info">
              Cetak Surat Perjanjian Kerja Sama
          </button>
        </td>
      </tr>`
      $('#tablePelanggan tbody').append(row);
    });
    // TODO listener button edit and delete
  });

  $('#search').on('keypress', function (e) {
    var value = $(this).val();
    if (e.keyCode == 13) {
      Daftar.Page = 1;
      Daftar.Keyword = value;
      Daftar.Apply();
      $(this).val('');
    }
  });

  $('#nextPage').on('click', function (e) {
    console.log(Daftar.Data.length, Daftar.Rows)
    if (Daftar.Data.length == Daftar.Rows) {
      Daftar.Page = Daftar.Page + 1;
      Daftar.Apply(function () {
        if (Daftar.Data.length == Daftar.Rows) {
          $('#nextPage').parent().removeClass('disabled');
        } else {
          $('#nextPage').parent().addClass('disabled');
        }
        if (Daftar.Page > 1) {
          $('#prevPage').parent().removeClass('disabled');
        } else {
          $('#prevPage').parent().addClass('disabled');
        }
      });
    }
  });

  $('#prevPage').on('click', function (e) {
    if (Daftar.Page > 1) {
      Daftar.Page = Daftar.Page - 1;
      Daftar.Apply(function () {
        if (Daftar.Data.length == Daftar.Rows) {
          $('#nextPage').parent().removeClass('disabled');
        } else {
          $('#nextPage').parent().addClass('disabled');
        }
        if (Daftar.Page > 1) {
          $('#prevPage').parent().removeClass('disabled');
        } else {
          $('#prevPage').parent().addClass('disabled');
        }
      });
    }
  });

  // TAMBAH
  $('#formAdd').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let tambah = $(this).serializeObject();
    tambah.project = Pelanggan.Form.project;
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
            $('#formAdd input').val('');
            $('#formAdd textarea').empty();
            $('#formAdd .projectRows').empty();
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

  var addProject = function (formHere, index, row) {
    let button = {
      tag: 'warning',
      text: 'Edit',
      class: 'edit-project',
      disabled: 'disabled',
    }
    if (!row) {
      row = {
        name: '',
        location: '',
      };
      button = {
        tag: 'success',
        text: 'OK',
        class: 'ok-project',
        disabled: '',
      };
    }
    console.log(index, row);
    let form = `<div class="form-row">
			<div class="form-group col-4">
				<label for="nameProject${index}">Nama</label>
				<input name="project[${index}][name]" type="text" class="form-control" id="nameProject${index}" value="${row.name}" required ${button.disabled}>
			</div>
			<div class="form-group col-6">
				<label for="locationProject${index}">Lokasi</label>
				<input name="project[${index}][location]" type="text" class="form-control" id="locationProject${index}" value="${row.name}" required ${button.disabled}>
			</div>
			<div class="col-2 d-flex align-items-end mb-3 btn-group">
        <button type="button" class="btn btn-${button.tag} btn-block ${button.class}" data-index="${index}">${button.text}</button>
        <button type="button" class="btn btn-danger btn-block delete-project" data-index="${index}">Delete</button>
			</div>
    </div>`;

    $(formHere).find('.projectRows').append(form);
    let editProject = function (here) {
      let row = $(this).parents('.form-row');
      $(row).find('input, select').removeAttr('disabled');
      $(row).find('.edit-project').removeClass('btn-warning');
      $(row).find('.edit-project').addClass('btn-success');
      $(row).find('.edit-project').html('OK');
      $(row).find('.edit-project').unbind('click');
      $(row).find('.edit-project').on('click', function () {
        okProject($(this));
      });
      $(row).find('.edit-project').addClass('ok-project');
      $(row).find('.edit-project').removeClass('edit-project');
    }
    let okProject = function (here) {
      let index = $(here).data('index');
      let row = $(here).parents('.form-row');
      $(row).find('input, select').attr('disabled', 'disabled');
      $(row).find('.ok-project').removeClass('btn-success');
      $(row).find('.ok-project').addClass('btn-warning');
      $(row).find('.ok-project').html('Edit');
      $(row).find('.ok-project').unbind('click');
      $(row).find('.ok-project').on('click', function () {
        editProject($(this));
      });
      $(row).find('.ok-project').addClass('edit-project');
      $(row).find('.ok-project').removeClass('ok-project');
      Pelanggan.EditProject(index, {
        name: $('#nameProject' + index).val(),
        location: $('#locationProject' + index).val(),
      });
    }
    let deleteProject = function (here) {
      $(here).parents('.form-row').remove();
      let index = $(here).data('index');
      Pelanggan.RemoveProject(index, function () {
        $(formHere).find('.projectRows').find('.form-row').each(function (i, row) {
          $(row).find('[id^="nameProject"').attr({
            id: `#nameProject${i}`,
            name: `price[rent][${i}][duration]`,
          });

          $(row).find('[id^="locationProject"').attr({
            id: `#locationProject${i}`,
            name: `price[rent][${i}][desc]`,
          });

          $(row).find('.ok-project, .delete-project, .edit-project').attr('data-index', `${i}`);
        });
      });
    }

    $('.edit-project').on('click', function () {
      editProject($(this));
    });
    $('.ok-project').on('click', function () {
      okProject($(this));
    });
    $('.delete-project').on('click', function () {
      deleteProject($(this));
    });
  }

  var setByType = function (form, type) {
    if (type == 1) {
      $(form).find('[name="group_name"]').parent().prop('hidden', true);
      $(form).find('[name="role"]').parent().prop('hidden', true);
      $(form).find('.projectButton').prop('disabled', true);
      $(form).find('.projectRows').html('');
    } else if (type == 2) {
      $(form).find('[name="group_name"]').parent().prop('hidden', false);
      $(form).find('[name="role"]').parent().prop('hidden', false);
      $(form).find('.projectButton').prop('disabled', false);
    }
  }

  $('#addProjectNew').on('click', function () {
    addProject($('#formAdd'), Pelanggan.LenProject());
    Pelanggan.AppendProject({
      name: '',
      location: '',
    });
  });

  $('#typeAdd').on('change', function () {
    let type = $(this).val();
    setByType($('#formAdd'), type);
  });

  // UBAH
  $('#formEdit').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let ubah = $(this).serializeObject();
    ubah.project = Pelanggan.Form.project;
    Pelanggan.Set(ubah);
    Pelanggan.Validate(true, (ok) => {
      Form.Reset($(this));
      if (ok.valid) {
        Loading.Start($('#formEdit [type="submit"]'));
        Pelanggan.Edit(() => {
          Form.Message('success', 'berhasil mengubah pelanggan', $('#messageEdit'));
          Loading.End();
          Form.Reset($('#formEdit'));
        }, () => {
          Form.Message('danger', 'gagal mengubah pelanggan', $('#messageEdit'));
          Loading.End();
        });
      } else {
        Form.Validate(ok.message, $('#formEdit'));
      }
    });
  });

  $('.autocomplete').autoComplete({
    resolverSettings: {
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });

  $('#typeEdit').on('change', function () {
    let type = $(this).val();
    setByType($('#formEdit'), type);
  });

  $('#formEdit .autocomplete').on('autocomplete.select', (e, cust) => {
    let name = $('#formEdit [name="name"]');
    let type = $('#formEdit [name="type"]');
    let groupName = $('#formEdit [name="group_name"]');
    let role = $('#formEdit [name="role"]');
    let address = $('#formEdit [name="address"]');
    let phone = $('#formEdit [name="phone"]');
    let nik = $('#formEdit [name="nik"]');
    let projects = $('#formEdit .projectRows');
    let submit = $('#formEdit [type="submit"]');
    if (!cust) {
      $(name).attr('disabled', 'disabled');
      $(name).val('');
      $(groupName).attr('disabled', 'disabled');
      $(groupName).val('');
      $(role).attr('disabled', 'disabled');
      $(role).val('');
      $(address).attr('disabled', 'disabled');
      $(address).val('');
      $(phone).attr('disabled', 'disabled');
      $(phone).val('');
      $(nik).attr('disabled', 'disabled');
      $(nik).val('');
      $(type).attr('disabled', 'disabled');
      $(submit).attr('disabled', 'disabled');
      return
    }
    Pelanggan.GetDetail(cust.value, (d) => {
      Pelanggan.Form = d
      console.log(d);
      Pelanggan.SetProject(addProject, function () {
        setByType($('#formEdit'), d.type);
        $(projects).empty();
      });

      $(name).val(d.name);
      $(type).val(d.type);
      $(groupName).val(d.group_name);
      $(role).val(d.role);
      $(nik).val(d.nik);
      $(address).val(d.address);
      $(phone).val(d.phone);
      $(name).removeAttr('disabled');
      $(type).removeAttr('disabled');
      $(groupName).removeAttr('disabled');
      $(role).removeAttr('disabled');
      $(nik).removeAttr('disabled');
      $(address).removeAttr('disabled');
      $(phone).removeAttr('disabled');
      $(submit).removeAttr('disabled');

    });
  });

});