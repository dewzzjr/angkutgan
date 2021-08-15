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
      if (!data.items || data.items.length == 0) {
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
    Bayar.SetFunc(Transactions.Get);
    Bayar.Init();
    Kirim.SetFunc(Transactions.Get);
    Kirim.Init();
  });

  // DAFTAR
  Daftar.Init('/rental');
  Daftar.GetData(function (data) {
    $('#tableTx tbody').empty();
    data.forEach(e => {
      let id = e.customer.code + '_' + e.tx_date.split('/').reverse().join('');
      let name = (e.customer.group_name) ? e.customer.group_name : e.customer.name;

      let total = e.total_price;
      let items = '';
      e.items.forEach(i => {
        items += `
        <li class="list-group-item text-right">
          ${i.name} &times; ${i.amount} ${i.item_unit} &times; ${i.duration} ${i.time_unit_desc} = ${formatPrice(i.amount * i.duration * i.price)}
        </li>`
      });
      if (e.shipping_fee) {
        total += e.shipping_fee;
        items += `
        <li class="list-group-item text-right">
          Ongkos Kirim = ${formatPrice(e.shipping_fee)}
        </li>
        `
      }
      if (e.deposit) {
        total += e.deposit;
        items += `
        <li class="list-group-item text-right">
          Deposit = ${formatPrice(e.deposit)}
        </li>
        `
      }
      let btnTx = (!e.status.in_payment) ? `
      <button type="button" class="btn btn-warning editBtn">Ubah</button>
      <button type="button" class="btn btn-danger deleteBtn">Hapus</button>
      ` : '';
      let btnPay = (!e.status.payment_done) ? `
      <button type="button" class="btn btn-success paymentBtn">Bayar</button>
      ` : '';
      let btnShip = (!e.status.shipping_done && e.status.in_payment && !e.status.in_shipping) ? `
      <button type="button" class="btn btn-primary shipmentBtn">Kirim</button>
      ` : (!e.status.shipping_done && e.status.in_payment) ? `
      <button type="button" class="btn btn-secondary shipmentBtn">Ubah Pengiriman</button>
      ` : '';
      let btnPrint = (e.status.payment_done) ? `
      <button type="button" class="btn btn-secondary print">Cetak</button>
      ` : '';
      let btnReturn = (e.status.in_shipping && !e.status.is_return) ? `
      <button type="button" class="btn btn-primary">Kembali</button>
      ` : '';
      let btnExtend = (e.status.in_shipping && !e.status.is_return && !e.status.done) ? `
      <button type="button" class="btn btn-warning" data-toggle="modal" data-target="#extendModal">
        Perpanjangan
      </button>` : '';
      let tx = `
			<tr class="rowTx" data-row="${id}">
        <td>${e.tx_date}</td>
        <td>${name}</td>
        <td>${e.status.desc}</td>
        <td>
          <div class="btn-group">
            <button type="button" class="btn btn-info" data-toggle="collapse" data-target="#${id}"
              data-parent="#tableTx" class="collapsed">
              Detail
            </button>
            ${btnTx}
          </div>
          <div class="btn-group">
            ${btnPay}${btnPrint}
          </div>
          <div class="btn-group">
            ${btnShip}${btnReturn}${btnExtend}
          </div>
        </td>
      </tr>
      <tr class="collapse" id="${id}">
        <td colspan="4">
          <div class="card">
            <div class="card-header">
              Ringkasan
            </div>
            <div class="card-body">
              <h5 class="card-title">${e.customer.code} - ${name}</h5>
              <h6 class="card-subtitle mb-2 text-muted">
                ${e.address}
              </h6>
              <p class="card-text">CP ${e.customer.name} - ${e.customer.phone}</p>
            </div>
            <ul class="list-group list-group-flush">
              ${items}
            </ul>
            <div class="card-footer text-right font-weight-bold">
              Total Tagihan: ${formatPrice(total)}
            </div>
          </div>
        </td>
      </tr>`
      $('#tableTx tbody').append(tx);
    });

    $('#tableTx .editBtn').on('click', function () {
      let code = $(this).closest('.rowTx').data('row').split('_');
      var query = {
        customer: code[0],
        date: code[1],
        action: 'create'
      };
      var url = window.location.pathname + '?' + $.param(query);
      window.location.replace(url);
    });

    $('#tableTx .paymentBtn').on('click', function (e) {
      let code = $(this).closest('.rowTx').data('row').split('_');
      var query = {
        customer: code[0],
        date: code[1],
        action: 'payment'
      };
      var url = window.location.pathname + '?' + $.param(query);
      window.location.replace(url);
    });

    $('#tableTx .shipmentBtn').on('click', function (e) {
      let code = $(this).closest('.rowTx').data('row').split('_');
      var query = {
        customer: code[0],
        date: code[1],
        action: 'shipment'
      };
      var url = window.location.pathname + '?' + $.param(query);
      window.location.replace(url);
    });

    $('#tableTx .print').on('click', function (e) {
      let code = $(this).closest('.rowTx').data('row').split('_');
      var query = {
        customer: code[0],
        date: code[1]
      };
      var template = '/KwitansiPembayaran.html?' + $.param(query);
      window.open(template);
    });

    $('#tableTx .deleteBtn').on('click', function (e) {
      let code = $(this).closest('.rowTx').data('row').split('_').join('/');
      $.confirm({
        title: 'Peringatan!',
        content: `Apakah anda yakin untuk menghapus "${code}"?`,
        buttons: {
          ok: function () {
            $.ajax({
              type: 'DELETE',
              url: `/rental/${code}`,
              contentType: 'application/json',
              success: function (data, status, xhr) {
                if (status === 'success' && data.result == 'OK') {
                  Daftar.Reload();
                  $.alert({
                    title: 'Berhasil',
                    content: `${code}: berhasil dihapus.`,
                  });
                }
                console.log(data, status);
              },
              error: function (xhr, status, error) {
                console.log(status, error);
                $.alert({
                  title: 'Gagal',
                  content: `${code}: ${error}`,
                });
              }
            });
          },
          cancel: function () {}
        },
      });
    });
  });

  var search = function (e) {
    let cust = $('#list [name="search_code"]').val();
    let date = $('#searchDate').val();
    console.log(cust, date);
    // TODO search
  }

  $('#searchDate').datepicker({
    format: 'dd/mm/yyyy',
  }).on('change', search);

  $('#searchCustomer').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });
  $('#searchCustomer').on('autocomplete.select', search);

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

  // SUBMIT
  serializeObject = () => {
    return {
      date: $('#create [name="tx_date"]').val(),
      customer: $('#create [name="customer"]').val(),
      address: $('#create [name="address"]').val(),
      project_id: parseInt($('#create [name="project_id"]').val()) || 0,
      discount: parseInt($('#create [name="discount"]').val()) || 0,
      shipping_fee: parseInt($('#create [name="shipping_fee"]').val()) || 0,
      deposit: parseInt($('#create [name="deposit"]').val()) || 0,
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
            clearTx();
            Form.Reset($('#formTx'));
          }, () => {
            Form.Message('danger', 'gagal mengubah transaksi', $('#messageTx'));
            Loading.End();
          });
        } else {
          Transactions.Create(() => {
            Form.Message('success', 'berhasil membuat transaksi', $('#messageTx'));
            Loading.End();
            clearTx();
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

  // DATE
  $('#tableBarang').hide();
  var getTransaction = function () {
    let customer = $('#create [name="customer"]').val();
    let date = $('#create [name="tx_date"]').val();
    if (customer && date) {
      Loading.Start($('#initialTx'));
      Transactions.Get(customer, date, (ok) => {
        setEdit(ok);
        if (ok.id) {
          Transactions.IsEdit = true;
        } else {
          Transactions.IsEdit = false;
        }
        Loading.End();
        $('#tableBarang').show();
      }, () => {
        Loading.End();
        $('#tableBarang').show();
      });
    } else {
      $('#tableBarang').hide();
    }
  }
  $('#create [name="tx_date"]').datepicker({
    format: 'dd/mm/yyyy',
  }).on('change', function () {
    getTransaction();
  });

  // CUSTOMER
  var initCustomer = function (code) {
    Transactions.GetCustomer(code, (customer) => {
      let name = (customer.group_name) ? customer.group_name : customer.name;
      $('.customerName').html(name);
      $('#customerName').val(name);
      if (customer.type == 2) {
        $('.contactPerson').html(`CP ${customer.name} (${customer.role}) - ${customer.phone}`)
      } else {
        $('.contactPerson').html(`${customer.phone}`);
      }

      $('#create [name="address"]').prop('disabled', false).val(customer.address);
      if (customer.project.length > 0) {
        $('#projectTx').parent().show();
        $('#create [name="project_id"]').html('<option value="" selected>Tidak ada proyek</option>');
        customer.project.forEach(e => {
          $('#create [name="project_id"]').append(`<option value="${e.id}">${e.name}</option>`);
        });
      } else {
        $('#projectTx').parent().hide()
      }
      summary();
    });
  }
  $('#customerCode').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });
  $('#customerCode').on('autocomplete.select', function (e) {
    getTransaction();
    let customer = $('#create [name="customer"]').val();
    if (customer) {
      initCustomer(customer);
    }
  });

  // INIT TRANSACTION
  if (Menu.Query['date'] && Menu.Query['customer']) {
    let date = Menu.Query['date'];
    date = [date.substring(0, 4), date.substring(4, 6), date.substring(6, 8)].join('-');
    $('#create [name="tx_date"]').datepicker('update', new Date(date));
    let cust = Menu.Query['customer'];
    $('#customerCode').autoComplete('set', {
      value: cust,
      text: cust
    });
    initCustomer(cust);
    getTransaction();
  } else {
    $('#create [name="tx_date"]').datepicker('update', new Date());
  }

  // EDIT
  var setEdit = (ok) => {
    if (ok.id) {
      if (ok.project_id) {
        $('#create [name="project_id"]').val(ok.project_id);
        $('#create [name="address"]').prop('disabled', true);
      } else {
        $('#create [name="project_id"]').val('');
        $('#create [name="address"]').prop('disabled', false);
      }
      $('#create [name="address"]').val(ok.address);
      $('#create [name="discount"]').val(ok.discount);
      $('#create [name="shipping_fee"]').val(ok.shipping_fee);
      $('#create [name="deposit"]').val(ok.deposit);
      Transactions.Items = ok.items;
    } else {
      Transactions.Items = [];
    }
    Transactions.Reload();
    summary();
  }

  // PROJECT
  $('#projectTx').parent().hide()
  $('#create [name="project_id"]').on('change', function (e) {
    let id = $(this).val();
    if (id == "") {
      $('#create [name="address"]').prop('disabled', false);
      $('#create [name="address"]').val(Transactions.Customer.address);
    } else {
      $('#create [name="address"]').prop('disabled', true);
      let project = Transactions.Customer.project.filter((e) => {
        return e.id == id;
      });
      if (project.length > 0) {
        $('#create [name="address"]').val(project[0].location);
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
      $('#create [name="duration"]').val(edit.duration);
      $('#create [name="amount"]').val(edit.amount);
      $('#itemAction').html(action);
      $('#itemSubmit').on('click', editItem);
      $('#itemCancel').on('click', clearItem);
    }
    Transactions.GetItem(code, (item) => {
      if (edit) {
        $('#addCodeBarang').autoComplete('set', {
          value: code,
          text: `${code} - ${item.name}`
        });
      }
      $('#itemUnit').html(item.unit);
      $('#create [name="name"]').val(item.name);
      $('#create [name="amount"]').attr('max', item.avail.inventory);
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
            $('#create [name="duration"]').val(e.duration);
            $('#create [name="time"]').val(newTime);
            $('#create [name="timeunit"]').val(e.unit);
            $('#create [name="unit"]').val(e.unit_desc);
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
        $('#create [name="duration"]').val(duration);
        $('#create [name="timeunit"]').val(timeunit);
        $('#create [name="unit"]').val(unit);
        calculateItem();
      });
      calculateItem();
    });
  }

  var getCalculation = function () {
    let price = parseInt($('#create [name="price"]:checked').val()) || 0;
    let time = parseInt($('#create [name="time"]').val()) || 0;
    let duration = parseInt($('#create [name="duration"]').val()) || 0;
    let amount = parseInt($('#create [name="amount"]').val()) || 0;
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
    let address = $('#create [name="address"]').val();
    $('.shippingLocation').html(address);
    // item
    $('.itemSummary').empty();
    let total = 0
    if (Transactions.Items) {
      Transactions.Items.forEach(e => {
        let price = e.price * e.duration * e.amount;
        total += price;
        let item = `
        <li class="list-group-item text-right">
          ${e.name} &times; ${e.amount} ${e.item_unit} &times; ${e.duration} ${e.time_unit_desc} = ${formatPrice(price)}
        </li>`
        $('.itemSummary').append(item);
      });
    }
    // discount
    let percent = parseInt($('#create [name="discount"]').val()) || 0;
    let discountPrice = total * percent * -1 / 100;
    let discount = `
    <li class="list-group-item text-right">
      Discount = ${formatPrice(discountPrice)}
    </li>`
    $('.itemSummary').append(discount);
    // shipping fee
    let shippingFee = parseInt($('#create [name="shipping_fee"]').val()) || 0;
    let ongkir = `
    <li class="list-group-item text-right">
      Ongkos Kirim = ${formatPrice(shippingFee)}
    </li>`
    $('.itemSummary').append(ongkir);
    // deposit
    let depositFee = parseInt($('#create [name="deposit"]').val()) || 0;
    let deposit = `
    <li class="list-group-item text-right">
      Deposit = ${formatPrice(depositFee)}
    </li>`
    $('.itemSummary').append(deposit);
    // total
    total += discountPrice + shippingFee + depositFee;
    $('#totalFee').html(formatPrice(total));
  }
  var clearItem = function () {
    $('#addCodeBarang').autoComplete('clear');
    $('#itemUnit').html('');
    $('#timeDuration').html('');
    $('.rentSelection').empty();
    $('#formItem input').val('');
    $('#create [name="amount"]').attr('max', 0);
    $('#itemAction').html('<button class="btn btn-success" id="itemSubmit" disabled>Tambah</input>');
    $('#itemSubmit').on('click', addItem);
  }
  var clearTx = function () {
    $('#customerCode').autoComplete('clear');
    $('#create [name="tx_date"]').datepicker('update', new Date());
    $('#formTx input, #formTx textarea, #formTx select').val('');
    $('#formTx input, #formTx textarea').val('');
    Transactions.Form = {};
    Transactions.Items = [];
    Transactions.Customer = {};
    Transactions.IsEdit = false;
    Transactions.Reload();
    clearItem();
    $('.itemSummary').empty();
    $('.item').remove();

    $('.customerName').empty();
    $('#customerName').val('');
    $('.contactPerson').empty();

    $('#create [name="address"]').val('').prop('disabled', true);
    $('#create [name="project_id"]').empty();
    $('#projectTx').parent().hide();

    $('.shippingLocation').empty();
    $('#totalFee').empty();
  }
  var calculateItem = function () {
    let i = getCalculation();
    $('#create [name="total"]').val(i.total);
    let code = $('#create [name="code"]').val();
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
    let timeunit = parseInt($('#create [name="timeunit"]').val()) || 0;
    let unitdesc = $('#create [name="unit"]').val();
    let itemName = $('#create [name="name"]').val();
    let code = $('#create [name="code"]').val();
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
    Transactions.AddItem(item, clearItem);
    summary();
  }
  var editItem = function (e) {
    e.preventDefault();
    let i = getCalculation();
    let timeunit = parseInt($('#create [name="timeunit"]').val()) || 0;
    let unitdesc = $('#create [name="unit"]').val();
    let itemName = $('#create [name="name"]').val();
    let code = $('#create [name="code"]').val();
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
    let index = parseInt($(this).data('index'));
    Transactions.EditItem(item, index, clearItem);
    summary();
  }
  $('#itemSubmit').on('click', addItem);
  Transactions.Reload = function () {
    $('.item').remove();
    if (!Transactions.Items) {
      return;
    }
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