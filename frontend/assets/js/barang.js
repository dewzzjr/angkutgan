const Barang = {
  Form: {},
  Set: function (data) {
    Barang.Form = data;
  },
  Clear: function (callback) {
    Barang.Form = {};
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
      if (!data.unit) {
        ok.valid = false;
        ok.message.push({
          name: "unit",
          text: "satuan tidak boleh kosong"
        });
      } else {
        data.unit = data.unit.toLowerCase();
      }
      if (!isInt(data.stock)) {
        ok.valid = false;
        ok.message.push({
          name: "stock",
          text: "stok harus berupa angka"
        });
      } else {
        data.stock = parseInt(data.stock);
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
          action: 'validate_code_item',
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
    let retries = false;
    $.ajax({
      type: 'POST',
      url: '/item',
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
  },
  Edit: function (callback, failedCallback) {
    let data = this.Form;
    let set = this.Set;
    let retries = false;
    $.ajax({
      type: 'PATCH',
      url: `/item/${data.code}`,
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
  },
  GetDetail: function (code, callback) {
    $.ajax({
      type: 'GET',
      url: `/item/${code}`,
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

const Harga = {
  Form: {
    code: 0,
    price: {
      sell: 0,
      rent: [],
    }
  },
  Set: function (form) {
    Harga.Form = form;
  },
  AppendRent: function (rent) {
    Harga.Form.price.rent.push(rent);
  },
  EditRent: function (index, rent) {
    Harga.Form.price.rent[index] = rent;
  },
  LenRent: function () {
    return this.Form.price.rent.length;
  },
  RemoveRent: function (index, callback) {
    let pop = Harga.Form.price.rent.splice(index, 1);
    if (callback) {
      callback(pop);
    }
  },
  SetRent: function (func, before, after) {
    if (before) {
      before();
    }
    Harga.Form.price.rent.forEach(function (r, i) {
      func(i, r);
    });
    if (after) {
      after();
    }
  },
  Validate: function (callback) {
    let ok = {
      message: [],
      valid: true
    };
    let data = this.Form;
    let check = function () {
      if (!isInt(data.price.sell)) {
        ok.valid = false;
        ok.message.push({
          name: "price[sell]",
          text: "harga jual harus berupa angka"
        });
      } else {
        data.price.sell = parseInt(data.price.sell);
      }

      (data.price.rent ? data.price.rent : []).forEach(function (r, i) {
        if (!r.desc) {
          ok.valid = false;
          ok.message.push({
            name: `price[rent][${i}][desc]`,
            text: "deskripsi tidak boleh kosong"
          });
        }

        if (!isInt(r.unit)) {
          ok.valid = false;
          ok.message.push({
            name: `price[rent][${i}][unit]`,
            text: "pilih satuan waktu sewa"
          });
        } else {
          data.price.rent[i].unit = parseInt(r.unit);
        }

        if (!isInt(r.value)) {
          ok.valid = false;
          ok.message.push({
            name: `price[rent][${i}][value]`,
            text: "harga sewa harus berupa angka"
          });
        } else {
          data.price.rent[i].value = parseInt(r.value);
        }

        if (!isInt(r.duration)) {
          ok.valid = false;
          ok.message.push({
            name: `price[rent][${i}][duration]`,
            text: "durasi sewa harus berupa angka"
          });
        } else {
          data.price.rent[i].duration = parseInt(r.duration);
        }
      });

      if (callback) {
        callback(ok);
      }
      console.log(data);
      Harga.Set(data);
    }

    check();
    return ok
  }
};

$(document).ready(function () {
  header();

  // DAFTAR
  Daftar.Init('/items');
  Daftar.GetData(function (data) {
    $('#tableBarang tbody').html('');
    data.forEach(e => {
      let rent = '';
      let rentLabel = ''
      if (e.price.rent) {
        let irent = ''
        e.price.rent.forEach((r) => {
          var format = formatPrice(r.value);
          let duration = '';
          if (r.duration > 1) {
            duration = `/ ${r.duration} ${r.unit_desc}`
          } else if (r.duration == 1) {
            duration = `/ ${r.unit_desc}`
          }
          irent += `<li>${r.desc} : ${format} ${duration}</li>`
        });
        if (irent) {
          rentLabel = `<p>Harga Sewa: </p>`
          rent = `<ul>${irent}</ul>`
        }
      }
      var format = formatPrice(e.price.sell);
      var row = `
      <tr>
        <th scope="row">${e.code}</th>
        <td>${e.name}</td>
        <td>${e.avail.inventory}/${e.avail.asset}</td>
        <td>${e.unit}</td>
        <td>
          <div class="btn-group">
            <button type="button" class="btn btn-info" data-toggle="collapse" data-target="#${e.code}"
              data-parent="#tableBarang" class="collapsed">Detail</button>
            <button type="button" class="btn btn-warning">Ubah</button>
            <button type="button" class="btn btn-danger">Hapus</button>
          </div>
        </td>
      </tr>
      <tr class="collapse" id="${e.code}">
        <th class="text-right" scope="row" colspan="2">
          <p>Harga Jual: </p>
          ${rentLabel}
        </th>
        <td colspan="3">
          <p>${format}</p>
          ${rent}
        </td>
      </tr>`
      $('#tableBarang tbody').append(row);
    });
    // TODO listener button edit and delete, cetak surat
  });

  $('#search').on('keypress', function(e) {
    var value = $(this).val();
    if (e.keyCode == 13) {
      Daftar.Page = 1;
      Daftar.Keyword = value;
      Daftar.Apply();
      $(this).val('');
    }
  });

  $('#nextPage').on('click', function(e) {
    console.log(Daftar.Data.length, Daftar.Rows)
    if (Daftar.Data.length == Daftar.Rows) {
      Daftar.Page = Daftar.Page + 1;
      Daftar.Apply(function() {
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

  $('#prevPage').on('click', function(e) {
    if (Daftar.Page > 1) {
      Daftar.Page = Daftar.Page - 1;
      Daftar.Apply(function() {
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
    console.log(tambah);
    Barang.Set(tambah);
    Barang.Validate(false, (ok) => {
      Form.Reset($('#formAdd'));
      if (ok.valid) {
        Loading.Start($('#formAdd [type="submit"]'));
        Barang.Create(() => {
          Form.Message('success', 'berhasil menambah barang', $('#messageAdd'));
          Loading.End();
          Form.Reset($('#formAdd'), () => {
            $('#formAdd input').val('');
          });
        }, () => {
          Form.Message('danger', 'gagal menambah barang', $('#messageAdd'));
          Loading.End();
        });
      } else {
        Form.Validate(ok.message, $('#formAdd'));
      }
    });
  });

  // UBAH
  $('.autocomplete').autoComplete({
    resolverSettings: {
      url: '/ajax?action=items',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });
  $('#formEdit .autocomplete').on('autocomplete.select', (e, item) => {
    let name = $('#formEdit [name="name"]');
    let unit = $('#formEdit [name="unit"]');
    let stock = $('#formEdit [name="stock"]');
    let submit = $('#formEdit [type=submit]');
    if (!item) {
      $(name).attr('disabled', 'disabled');
      $(unit).attr('disabled', 'disabled');
      $(stock).attr('disabled', 'disabled');
      $(submit).attr('disabled', 'disabled');
      $(name).val('');
      $(unit).val('');
      $(stock).val('');
      return
    }
    Barang.GetDetail(item.value, (d) => {
      $(name).val(d.name);
      $(name).removeAttr('disabled');

      $(unit).val(d.unit);
      $(unit).removeAttr('disabled');

      let disewa = d.avail.asset - d.avail.inventory;
      $(stock).val(d.avail.asset);
      $(stock).attr('min', disewa);
      $(stock).removeAttr('disabled');
      $(stock).parent().find('.form-text').html(`Jumlah barang yang sekarang sedang disewa adalah ${disewa} ${d.unit}`);

      $(submit).removeAttr('disabled');
    });
  });
  $('#formEdit').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let ubah = $(this).serializeObject();
    Barang.Set(ubah);
    Barang.Validate(true, (ok) => {
      Form.Reset($(this));
      if (ok.valid) {
        Loading.Start($('#formEdit [type="submit"]'));
        Barang.Edit(() => {
          Form.Message('success', 'berhasil mengubah barang', $('#messageEdit'));
          Loading.End();
          Form.Reset($('#formEdit'));
        }, () => {
          Form.Message('danger', 'gagal mengubah barang', $('#messageEdit'));
          Loading.End();
        });
      } else {
        Form.Validate(ok.message, $('#formEdit'));
      }
    });
  });

  // HARGA
  var addRent = function (index, row) {
    let button = {
      tag: 'warning',
      text: 'Edit',
      class: 'edit-rent',
      disabled: 'disabled',
    }
    if (!row) {
      row = {
        desc: '',
        duration: 0,
        unit: 1,
        value: 0,
      };
      button = {
        tag: 'success',
        text: 'OK',
        class: 'ok-rent',
        disabled: '',
      };
    }
    let form = `<div class="form-row">
			<div class="form-group col-4">
				<label for="descRent">Deskripsi</label>
				<input name="price[rent][${index}][desc]" type="text" class="form-control" id="descRent${index}" value="${row.desc}" required ${button.disabled}>
			</div>
			<div class="form-group col-3">
				<label for="timeRent">Waktu</label>
				<div class="input-group">
					<input name="price[rent][${index}][duration]" type="number" class="form-control" id="timeRent${index}" value="${row.duration}" ${button.disabled}>
					<select name="price[rent][${index}][unit]" class="form-control" id="unitRent${index}" ${button.disabled}>
						<option value="1" ${row.unit == 1 ? "selected" : ""}>Minggu</option>
						<option value="2" ${row.unit == 2 ? "selected" : ""}>Bulan</option>
					</select>
				</div>
			</div>
			<div class="form-group col-3">
				<label for="priceRent">Harga Sewa</label>
				<div class="input-group">
					<div class="input-group-prepend">
						<span class="input-group-text">Rp</span>
					</div>
					<input name="price[rent][${index}][value]" type="number" class="form-control" id="priceRent${index}" value="${row.value}" ${button.disabled}>
				</div>
			</div>
			<div class="col-2 d-flex align-items-end mb-3 btn-group">
        <button type="button" class="btn btn-${button.tag} btn-block ${button.class}" data-index="${index}">${button.text}</button>
        <button type="button" class="btn btn-danger btn-block delete-rent" data-index="${index}">Delete</button>
			</div>
    </div>`;

    $('#rentRows').append(form);    
    let editRent = function (here) {
      let row = $(this).parents('.form-row');
      $(row).find('input, select').removeAttr('disabled');
      $(row).find('.edit-rent').removeClass('btn-warning');
      $(row).find('.edit-rent').addClass('btn-success');
      $(row).find('.edit-rent').html('OK');
      $(row).find('.edit-rent').unbind('click');
      $(row).find('.edit-rent').on('click', function() {
        okRent($(this));
      });
      $(row).find('.edit-rent').addClass('ok-rent');
      $(row).find('.edit-rent').removeClass('edit-rent');
    }
    let okRent = function (here) {
      let index = $(here).data('index');
      let row = $(here).parents('.form-row');
      $(row).find('input, select').attr('disabled', 'disabled');
      $(row).find('.ok-rent').removeClass('btn-success');
      $(row).find('.ok-rent').addClass('btn-warning');
      $(row).find('.ok-rent').html('Edit');
      $(row).find('.ok-rent').unbind('click');
      $(row).find('.ok-rent').on('click', function() {
        editRent($(this));
      });
      $(row).find('.ok-rent').addClass('edit-rent');
      $(row).find('.ok-rent').removeClass('ok-rent');
      Harga.EditRent(index, {
        desc: $('#descRent' + index).val(),
        unit: $('#unitRent' + index).val(),
        duration: $('#timeRent' + index).val(),
        value: $('#priceRent' + index).val(),
      });
    }
    let deleteRent = function (here) {
      $(here).parents('.form-row').remove();
      let index = $(here).data('index');
      Harga.RemoveRent(index, function () {
        $('#rentRows').find('.form-row').each(function (i, row) {
          $(row).find('[id^="timeRent"').attr({
            id: `#timeRent${i}`,
            name: `price[rent][${i}][duration]`,
          });

          $(row).find('[id^="descRent"').attr({
            id: `#descRent${i}`,
            name: `price[rent][${i}][desc]`,
          });

          $(row).find('[id^="unitRent"').attr({
            id: `#unitRent${i}`,
            name: `price[rent][${i}][unit]`,
          });

          $(row).find('[id^="priceRent"').attr({
            id: `#priceRent${i}`,
            name: `price[rent][${i}][value]`,
          });

          $(row).find('.ok-rent, .delete-rent, .edit-rent').attr('data-index', `${i}`);
        });
      });
    }

    $('.edit-rent').on('click', function() {
      editRent($(this));
    });
    $('.ok-rent').on('click', function() {
      okRent($(this));
    });
    $('.delete-rent').on('click', function () {
      deleteRent($(this));
    });
  }
  $('#addRent').on('click', function () {
    addRent(Harga.LenRent());
    Harga.AppendRent({
      desc: '',
      unit: '1',
      duration: '0',
      value: '0',
    });
  });
  $('#formPrice .autocomplete').on('autocomplete.select', (e, item) => {
    let name = $('#namePrice');
    let unit = $('#unitPrice');
    let sell = $('#formPrice [name="price[sell]"]');
    let btn = $('#formPrice [type=submit], #formPrice [type=button]');
    if (!item) {
      $(btn).attr('disabled', 'disabled');
      $(name).val('');
      $(unit).val('');
      $(sell).val('').attr('disabled', 'disabled');
      $('#rentRows').empty();
      return
    }
    Barang.GetDetail(item.value, (d) => {
      Harga.Form = {
        code: d.code,
        price: d.price,
      };

      Harga.SetRent(addRent, function() {
        $('#rentRows').empty();
      });
      $(name).val(d.name);
      $(unit).val(d.unit);
      $(sell).val(d.price.sell).removeAttr('disabled');
      $(btn).removeAttr('disabled');
    });
  });
  $('#formPrice').on('submit', function (e) {
    e.preventDefault();
    Form.Init($(this));
    Form.Reset();
    let harga = $(this).serializeObject();
    harga.price.rent = Harga.Form.price.rent;
    Harga.Set(harga);
    Harga.Validate((ok) => {
      Form.Reset();
      if (ok.valid) {
        Loading.Start($('#formPrice [type="submit"]'));
        Barang.Set(Harga.Form);
        Barang.Edit(() => {
          Form.Message('success', 'berhasil mengubah harga barang');
          Loading.End();
          Form.Reset();
        }, () => {
          Form.Message('danger', 'gagal mengubah harga barang');
          Loading.End();
        });
      } else {
        Form.Validate(ok.message);
      }
    });
  });

  // TAB
  // let code = Menu.Query['code'];
  // if (Menu.Query['action'] == 'edit') {
  //   if (code) {
  //     console.log(code);
  //     $('#formEdit .autocomplete').autoComplete('set', {
  //       value: code,
  //       text: code,
  //     });
  //     $('#formEdit .autocomplete').trigger('autocomplete.select');
  //   }
  // }
});