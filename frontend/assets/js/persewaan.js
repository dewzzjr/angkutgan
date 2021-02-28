const Transactions = {
  Customer: {},
  GetCustomer: function (code, callback, failedCallback) {
    $.ajax({
      type: 'GET',
      url: `/customer/${code}`,
      contentType: 'application/json',
      success: function (data, status, xhr) {
        if (status === 'success') {
          Transactions.Customer = data.result;
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
  Items: [],
  GetItem: function (code, callback, failedCallback) {
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
  },
  Reload: () => {},
  AddItem: function (data, callback) {
    if (
      data.code == '' ||
      data.price == 0 ||
      data.duration == 0 ||
      data.amount == 0 ||
      data.time_unit == 0
    ) {
      return
    }
    Transactions.Items.push(data);
    Transactions.Reload();
    if (callback) {
      callback();
    }
  },
  RemoveItem: function (index) {
    let pop = Transactions.Items.splice(index, 1);
    Transactions.Reload();
    return pop;
  },
  EditItem: function (data, index, callback) {
    if (
      data.code == '' ||
      data.price == 0 ||
      data.duration == 0 ||
      data.amount == 0 ||
      data.time_unit == 0
    ) {
      console.log(data);
      return
    }
    let pop = Transactions.Items.splice(index, 1, data);
    Transactions.Reload();
    if (callback) {
      callback();
    }
    console.log(pop);
    return pop;
  },
  IsEdit: false,
  Form: {},
  Create: (callback, failedCallback) => {
    let data = Transactions.Form;
    $.ajax({
      type: 'POST',
      url: '/rental',
      contentType: 'application/json',
      data: JSON.stringify(data),
      success: function (data, status, xhr) {
        if (status === 'success') {
          Transactions.Form = data.result;
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
          Auth.Refresh(function () {
            $.ajax(this);
          });
        }
      },
    });
  },
  Edit: (callback, failedCallback) => {
    let data = Transactions.Form;
    $.ajax({
      type: 'PATCH',
      url: '/rental',
      contentType: 'application/json',
      data: JSON.stringify(data),
      success: function (data, status, xhr) {
        if (status === 'success') {
          Transactions.Form = data.result;
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
          Auth.Refresh(function () {
            $.ajax(this);
          });
        }
      },
    });
  },
  Get: (customer, date, callback, failedCallback) => {
    let format = date.split("/");
    date = format[2] + format[1] + format[0];
    console.log(date);
    $.ajax({
      type: 'GET',
      url: `/rental/${customer}/${date}`,
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
  },
  Validate: function (callback) {
    let ok = {
      message: [],
      valid: true
    };
    let data = this.Form
    let check = () => {
      if (!data.date) {
        ok.valid = false;
        ok.message.push({
          name: "tx_date",
          text: "tanggal transaksi tidak boleh kosong"
        });
      } else if (data.date.split("/").length != 3) {
        ok.valid = false;
        ok.message.push({
          name: "tx_date",
          text: "format tanggal salah"
        });
      }
      if (!data.customer) {
        ok.valid = false;
        ok.message.push({
          name: "customer",
          text: "pelanggan tidak boleh kosong"
        });
      }
      if (!data.address) {
        ok.valid = false;
        ok.message.push({
          name: "address",
          text: "alamat tidak boleh kosong"
        });
      }
      if (!isInt(data.project_id)) {
        ok.valid = false;
        ok.message.push({
          name: "project_id",
          text: "proyek tidak valid"
        });
      }
      if (!isInt(data.discount) || data.discount > 100 || data.discount < 0) {
        ok.valid = false;
        ok.message.push({
          name: "discount",
          text: "discount tidak valid"
        });
      }
      if (!isInt(data.deposit) || data.deposit < 0) {
        ok.valid = false;
        ok.message.push({
          name: "deposit",
          text: "deposit tidak valid"
        });
      }
      if (!isInt(data.shipping_fee) || data.shipping_fee < 0) {
        ok.valid = false;
        ok.message.push({
          name: "shipping_fee",
          text: "ongkos kirim tidak valid"
        });
      }
      if (data.items.length == 0) {
        ok.valid = false;
        ok.message.push({
          name: "code_text",
          text: "barang belum dimasukkan"
        });
      }
      if (callback) {
        callback(ok);
      }
      Transactions.Form = data;
    }
    check();
    return ok;
  }
};
$(document).ready(function () {
  header(() => {
    $('#formTx :input').on('change', summary);
    $('#formItem :input').on('change', calculateItem);
  });
  // SUBMIT
  serializeObject = () => {
    return {
      date: $('[name="tx_date"]').val(),
      customer: $('[name="customer"]').val(),
      address: $('[name="address"]').val(),
      project_id: parseInt($('[name="project_id"]').val()) || 0,
      discount: parseInt($('[name="discount"]').val()) || 0,
      shipping_fee: parseInt($('[name="shipping_fee"]').val()) || 0,
      deposit: parseInt($('[name="deposit"]').val()) || 0,
      items: Transactions.Items,
    }
  }
  $('#formTx').on('submit', function (e) {
    e.preventDefault();
    Form.Reset($(this));
    let data = serializeObject();
    Transactions.Form = data;
    console.log(data);
    Transactions.Validate((ok) => {
      console.log(ok);
      Form.Reset($(this));
      if (ok.valid) {
        Loading.Start($('#formTx [type="submit"]'));
        if (Transactions.IsEdit) {
          Transactions.Edit(() => {
            Form.Message('success', 'berhasil mengubah transaksi', $('#messageTx'));
            Loading.End();
            Form.Reset($('#formTx'));
          }, () => {
            Form.Message('danger', 'gagal mengubah transaksi', $('#messageTx'));
            Loading.End();
          });
        } else {
          Transactions.Create(() => {
            Form.Message('success', 'berhasil membuat transaksi', $('#messageTx'));
            Loading.End();
            Form.Reset($('#formTx'));
          }, () => {
            Form.Message('danger', 'gagal membuat transaksi', $('#messageTx'));
            Loading.End();
          });
        }
      } else {
        Form.Validate(ok.message, $('#formTx'));
      }
    });
  });

  // CUSTOMER
  $('#customerCode').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });
  $('#customerCode').on('autocomplete.select', function (e, select) {
    if (!select) {
      return;
    }
    Transactions.GetCustomer(select.value, (customer) => {
      getTransaction(e);
      let name = (customer.group_name) ? customer.group_name : customer.name;
      $('.customerName').html(name);
      $('#customerName').val(name);
      if (customer.type == 2) {
        $('.contactPerson').html(`CP ${customer.name} (${customer.role}) - ${customer.phone}`)
      } else {
        $('.contactPerson').html(`${customer.phone}`);
      }

      $('[name="address"]').prop('disabled', false).val(customer.address);
      if (customer.project.length > 0) {
        $('#projectTx').parent().show();
        $('[name="project_id"]').html('<option value="" selected>Tidak ada proyek</option>');
        customer.project.forEach(e => {
          $('[name="project_id"]').append(`<option value="${e.id}">${e.name}</option>`);
        });
      } else {
        $('#projectTx').parent().hide()
      }
      summary();
    });
  });

  // EDIT
  var setEdit = (ok) => {
    console.log(ok);
    if (ok.project_id) {
      $('[name="project_id"]').val(ok.project_id);
      $('[name="address"]').prop('disabled', true);
    } else {
      $('[name="project_id"]').val('');
      $('[name="address"]').prop('disabled', false);
    }
    $('[name="address"]').val(ok.address);
    $('[name="discount"]').val(ok.discount);
    $('[name="shipping_fee"]').val(ok.shipping_fee);
    $('[name="deposit"]').val(ok.deposit);
    Transactions.Items = ok.items;
    Transactions.Reload();
    summary();
  }

  // DATE
  $('#formTx').hide();
  var getTransaction = (e) => {
    let customer = $('[name="customer"]').val();
    let date = $('[name="tx_date"]').val();
    if (customer && date) {
      Loading.Start($('#initialTx'));
      Transactions.Get(customer, date, (ok) => {
        setEdit(ok);
        Transactions.IsEdit = true;
        Loading.End();
        $('#formTx').show();
      }, () => {
        Transactions.IsEdit = false;
        Loading.End();
        $('#formTx').show();
      });
    } else {
      $('#formTx').hide();
    }
  }
  $('[name="tx_date"]').datepicker({
    format: 'dd/mm/yyyy',
  }).on('change', getTransaction);
  $('[name="tx_date"]').datepicker('update', new Date());

  // PROJECT
  $('#projectTx').parent().hide()
  $('[name="project_id"]').on('change', function (e) {
    let id = $(this).val();
    if (id == "") {
      $('[name="address"]').prop('disabled', false);
      $('[name="address"]').val(Transactions.Customer.address);
    } else {
      $('[name="address"]').prop('disabled', true);
      let project = Transactions.Customer.project.filter((e) => {
        return e.id == id;
      });
      if (project.length > 0) {
        $('[name="address"]').val(project[0].location);
      }
    }
    summary();
  });

  // ITEMS
  $('#addCodeBarang').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=items',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });
  $('#addCodeBarang').on('autocomplete.select', function (ev, select) {
    if (!select) {
      return;
    }
    initItem(select.value);
    summary();
  });
  var initItem = function (code, edit) {
    if (edit) {
      let index = edit.index;
      let action = `
      <button class="btn btn-warning" id="itemSubmit" disabled data-index="${index}">Ubah</input>
      <button class="btn btn-danger" id="itemCancel">Batal</input>
      `
      $('[name="duration"]').val(edit.duration);
      $('[name="amount"]').val(edit.amount);
      $('#itemAction').html(action);
      $('#itemSubmit').on('click', editItem);
      $('#itemCancel').on('click', clear);
    }
    Transactions.GetItem(code, (item) => {
      if (edit) {
        $('#addCodeBarang').autoComplete('set', {
          value: code,
          text: `${code} - ${item.name}`
        });
      }
      $('#itemUnit').html(item.unit);
      $('[name="name"]').val(item.name);
      $('[name="amount"]').attr('max', item.avail.inventory);
      $('.rentSelection').empty();
      let ok = false;
      item.price.rent.forEach(e => {
        let price = formatPrice(e.value);
        let id = (e.unit == 1) ? 'week' + e.duration : (e.unit == 2) ? 'month' + e.duration : 'other' + e.duration;
        if (edit) {
          let cond1 = edit.time_unit == e.unit;
          let cond2 = !ok;
          if (cond1 && cond2) {
            ok = true;
            let newTime = edit.duration / e.duration;
            $('#timeDuration').html(`&times; ${e.duration} ${e.unit_desc}`);
            $('[name="duration"]').val(e.duration);
            $('[name="time"]').val(newTime);
            $('[name="timeunit"]').val(e.unit);
            $('[name="unit"]').val(e.unit_desc);
          }
        }
        let checked = (ok) ? 'checked' : '';
        let rent = `
        <div class="form-check m-2 rentOption"
          data-duration="${e.duration}" 
          data-timeunit="${e.unit}" 
          data-unit="${e.unit_desc}">
          <input class="form-check-input" type="radio" name="price" 
            id="${id}" value="${e.value}" ${checked}>
          <label class="form-check-label" for="${id}">
            ${e.desc}<br />${e.duration} ${e.unit_desc} - ${price}
          </label>
        </div>`;
        $('.rentSelection').append(rent);
      });
      $('.rentOption').on('change', function (e) {
        let timeunit = $(this).data('timeunit');
        let unit = $(this).data('unit');
        let duration = $(this).data('duration');
        $('#timeDuration').html(`&times; ${duration} ${unit}`);
        $('[name="duration"]').val(duration);
        $('[name="timeunit"]').val(timeunit);
        $('[name="unit"]').val(unit);
        calculateItem();
      });
      calculateItem();
    });
  }

  var getCalculation = function () {
    let price = parseInt($('[name="price"]:checked').val()) || 0;
    let time = parseInt($('[name="time"]').val()) || 0;
    let duration = parseInt($('[name="duration"]').val()) || 0;
    let amount = parseInt($('[name="amount"]').val()) || 0;
    let total = price * amount * time;
    return {
      price: price,
      time: time,
      duration: duration,
      amount: amount,
      total: total,
    };
  }
  var summary = function (e) {
    // address
    let address = $('[name="address"]').val();
    $('.shippingLocation').html(address);
    // item
    $('.itemSummary').empty();
    let total = 0
    Transactions.Items.forEach(e => {
      let price = e.price * e.duration * e.amount;
      total += price;
      let item = `
      <li class="list-group-item text-right">
        ${e.name} &times; ${e.amount} ${e.item_unit} &times; ${e.duration} ${e.time_unit_desc} = ${formatPrice(price)}
      </li>`
      $('.itemSummary').append(item);
    });
    // discount
    let percent = parseInt($('[name="discount"]').val()) || 0;
    let discountPrice = total * percent * -1 / 100;
    let discount = `
    <li class="list-group-item text-right">
      Discount = ${formatPrice(discountPrice)}
    </li>`
    $('.itemSummary').append(discount);
    // shipping fee
    let shippingFee = parseInt($('[name="shipping_fee"]').val()) || 0;
    let ongkir = `
    <li class="list-group-item text-right">
      Ongkos Kirim = ${formatPrice(shippingFee)}
    </li>`
    $('.itemSummary').append(ongkir);
    // deposit
    let depositFee = parseInt($('[name="deposit"]').val()) || 0;
    let deposit = `
    <li class="list-group-item text-right">
      Deposit = ${formatPrice(depositFee)}
    </li>`
    $('.itemSummary').append(deposit);
    // total
    total += discountPrice + shippingFee + depositFee;
    $('#totalFee').html(formatPrice(total));
  }
  var clear = function () {
    $('#addCodeBarang').autoComplete('clear');
    $('#itemUnit').html('');
    $('#timeDuration').html('');
    $('.rentSelection').empty();
    $('#formItem input').val('');
    $('[name="amount"]').attr('max', 0);
    $('#itemAction').html('<button class="btn btn-success" id="itemSubmit" disabled>Tambah</input>');
    $('#itemSubmit').on('click', addItem);
  }
  var calculateItem = function () {
    let i = getCalculation();
    $('[name="total"]').val(i.total);
    let code = $('[name="code"]').val();
    if (i.total > 0 && code != '') {
      $('#itemSubmit').prop('disabled', false);
    } else {
      $('#itemSubmit').prop('disabled', true);
    }
    summary();
  }
  var addItem = function (e) {
    e.preventDefault();
    let i = getCalculation();
    let timeunit = parseInt($('[name="timeunit"]').val()) || 0;
    let unitdesc = $('[name="unit"]').val();
    let itemName = $('[name="name"]').val();
    let code = $('[name="code"]').val();
    let itemunit = $('#itemUnit').html();
    let item = {
      code: code,
      name: itemName,
      price: i.price / i.duration,
      amount: i.amount,
      item_unit: itemunit,
      time_unit: timeunit,
      time_unit_desc: unitdesc,
      duration: i.time * i.duration,
    };
    Transactions.AddItem(item, clear);
    summary();
  }
  var editItem = function (e) {
    e.preventDefault();
    let i = getCalculation();
    let timeunit = parseInt($('[name="timeunit"]').val()) || 0;
    let unitdesc = $('[name="unit"]').val()
    let code = $('[name="code"]').val();
    let itemunit = $('#itemUnit').html();
    let item = {
      code: code,
      price: i.price / i.duration,
      amount: i.amount,
      item_unit: itemunit,
      time_unit: timeunit,
      time_unit_desc: unitdesc,
      duration: i.time * i.duration,
    };
    let index = parseInt($(this).data('index'));
    Transactions.EditItem(item, index, clear);
    summary();
  }
  $('#itemSubmit').on('click', addItem);
  Transactions.Reload = function () {
    $('.item').remove();
    let len = Transactions.Items.length;
    Transactions.Items.slice().reverse().forEach((e, i) => {
      let index = len - i - 1;
      let price = formatPrice(e.price);
      let total = formatPrice(e.price * e.duration * e.amount);
      let item = `
			<tr class="item">
        <th scope="row">${e.code}</th>
        <td>${e.amount} ${e.item_unit}</td>
        <td>${e.duration} ${e.time_unit_desc}</td>
        <td>${price}/${e.time_unit_desc}</td>
        <th>${total}</th>
        <td>
          <button type="button" class="btn btn-warning editItemBtn" data-index="${index}">Ubah</button>
          <button type="button" class="btn btn-danger deleteItemBtn" data-index="${index}">Hapus</button>
        </td>
      </tr>`
      $('.items').prepend(item);
    });
    $('.editItemBtn').on('click', function () {
      let index = parseInt($(this).data('index'));
      let edit = Transactions.Items[index];
      edit.index = index;
      initItem(edit.code, edit);
    });
    $('.deleteItemBtn').on('click', function () {
      let index = parseInt($(this).data('index')) || -1;
      let row = $(this).parents().find('.item');
      let code = $(row).find('th[scope="row"]').html();
      $.confirm({
        title: 'Peringatan!',
        content: `Apakah anda yakin untuk menghapus "${code}"?`,
        buttons: {
          ok: function () {
            Transactions.RemoveItem(index);
            summary();
          },
          cancel: function () {}
        },
      });
    });
  }
});