const Bayar = {
  Form: {},
  Tx: {},
  SetFunc: (func) => {
    Bayar.GetTx = func;
  },
  Create: (callback, failedCallback) => {
    if (!Bayar.Tx.id) {
      return
    }
    let id = Bayar.Tx.id;
    let data = Bayar.Form;
    $.ajax({
      type: 'POST',
      url: `/payment/${id}`,
      contentType: 'application/json',
      data: JSON.stringify(data),
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
  Edit: (callback, failedCallback) => {
    if (!Bayar.Tx.id) {
      return
    }
    let id = Bayar.Tx.id;
    let data = Bayar.Form;
    $.ajax({
      type: 'PATCH',
      url: `/payment/${id}`,
      contentType: 'application/json',
      data: JSON.stringify(data),
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
  Reload: () => {
    let format = Bayar.Tx.tx_date.split("/");
    date = format[2] + format[1] + format[0];
    var query = {
      customer: Bayar.Tx.customer.code,
      date: date,
      action: 'payment'
    };
    var url = window.location.pathname + '?' + $.param(query);
    window.location.replace(url);
  }
}

Bayar.Init = function() {
  var getTransaction = function() {
    let customer = $('#payment [name="customer"]').val();
    let date = $('#payment [name="tx_date"]').val();
    if (customer && date) {
      Bayar.GetTx(customer, date, (ok) => {
        if (!ok.id) {
          return
        }
        $('#tableBayar tbody').empty();
        let total = ok.deposit + ok.shipping_fee + ok.total_price
        let paid = ok.payment.reduce((t, p) => {
          let amount = p.amount
          if (p.account == 100) {
            amount = -amount
          }
          return t + amount
        }, 0);
        ok.unpaid = total - paid;
        $('#payment [name="unpaid"]').val(formatPrice(ok.unpaid));
        let name = ok.customer.group_name ? ok.customer.group_name : ok.customer.name;
        $('#customerNamePayment').val(name);
        $('#paymentModal .modal-title').html(ok.tx_date +" "+ name);
        $('#formPayment [name="amount"]').attr('max', ok.unpaid);
        if (ok.unpaid <= 0) {
          $('#payment .addPayment').prop('disabled', true);
        } else {
          $('#payment .addPayment').prop('disabled', false);
        }
        ok.payment.forEach((e, idx, arr) => {
          let last = (arr.length - 1) == idx;

          // <button type="button" class="btn btn-warning editPayment" data-toggle="modal" data-target="#paymentModal">Ubah</button>
          let acceptBy = last ? `
          <div class="btn-group">
            <button type="button" class="btn btn-danger deletePayment" data-id="${ok.id}">Hapus</button>
          </div>
          ` : e.accept_by;
          let row = `
          <tr>
            <th scope="row">${e.date}</th>
            <td>${e.name}</td>
            <td>${e.method_desc}</td>
            <td class="text-right">${formatPrice(e.amount)}</td>
            <td>
              ${acceptBy}
            </td>
          </tr>`
          $('#tableBayar tbody').append(row);
        });
        // $('#tableBayar .editPayment').on('click', function() {});
        $('#tableBayar .deletePayment').on('click', function() {
          let id = $(this).data('id');
          $.confirm({
            title: 'Peringatan!',
            content: `Apakah anda yakin untuk menghapus pembayaran terakhir di transaksi ini?`,
            buttons: {
              ok: function () {
                $.ajax({
                  type: 'DELETE',
                  url: `/payment/${id}`,
                  contentType: 'application/json',
                  success: function (data, status, xhr) {
                    if (status === 'success' && data.result == 'OK') {
                      $.alert({
                        title: 'Berhasil',
                        content: `berhasil menghapus pembayaran`,
                        onClose: function() {
                          Bayar.Reload();
                        },
                      });
                    }
                    console.log(data, status);
                  },
                  error: function (xhr, status, error) {
                    console.log(status, error);
                    $.alert({
                      title: 'Gagal',
                      content: error,
                    });
                  }
                });
              },
              cancel: function () {}
            },
          });
        });
        Bayar.Tx = ok
      });
    }
  }

  $('#payment [name="tx_date"]').datepicker({
    format: 'dd/mm/yyyy',
  }).on('change', function() {
    getTransaction();
  });

  $('#customerCodePayment').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  }).on('autocomplete.select', function() {
    getTransaction();
  });

  if (Menu.Query['date'] && Menu.Query['customer'] && Menu.Query["action"] == 'payment') {
    let date = Menu.Query['date'];
    date = [date.substring(0,4),date.substring(4,6),date.substring(6,8)].join('-');
    $('#payment [name="tx_date"]').datepicker('update', new Date(date));
    let cust = Menu.Query['customer'];
    $('#customerCodePayment').autoComplete('set', { value: cust, text: cust });
    getTransaction();
  } else {
    $('#payment [name="tx_date"]').datepicker('update', new Date());
  }

  $('#payment .addPayment').on('click', () => {
    $('#paymentModal [name="date"]').datepicker({ format: 'dd/mm/yyyy' });
    $('#paymentModal [type="submit"]').data('action', 'add');
  });
  $('#payment .editPayment').on('click', () => {
    $('#paymentModal [name="date"]').datepicker({ format: 'dd/mm/yyyy' });
    $('#paymentModal [type="submit"]').data('action', 'edit');
  });

  var serializeObject = function () {
    return {
      date: $('#formPayment [name="date"]').val(),
      name: $('#formPayment [name="name"]').val(),
      method: parseInt($('#formPayment [name="method"]').val()) || 100,
      amount: parseInt($('#formPayment [name="amount"]').val()) || 0,
      account: parseInt($('#formPayment [name="account"]').val()) || 0,
    }
  }
  $('#formPayment').on('submit', (e) => {
    e.preventDefault();
    let action = $('#paymentModal [type="submit"]').data('action');
    Bayar.Form = serializeObject();
    Form.Reset('#paymentIndex');
    $('#paymentModal').modal('hide');
    if (action == 'edit') {
      Bayar.Edit(() => {
        Form.Message('success', 'berhasil mengubah pembayaran', $('#messagePay'));
        Form.Reset($('#formPayment'));
        $.alert({
          title: 'Berhasil',
          content: `berhasil mengubah pembayaran`,
          onClose: function() {
            Bayar.Reload();
          },
        });
        Bayar.Reload();
      }, () => {
        Form.Message('danger', 'gagal mengubah pembayaran', $('#messagePay'));
        Loading.End();
      });
    } else if (action == 'add') {
      Bayar.Create(() => {
        Form.Message('success', 'berhasil menambah pembayaran', $('#messagePay'));
        Form.Reset($('#formPayment'));
        $.alert({
          title: 'Berhasil',
          content: `berhasil menambah pembayaran`,
          onClose: function() {
            Bayar.Reload();
          },
        });
      }, () => {
        Form.Message('danger', 'gagal menambah pembayaran', $('#messagePay'));
      });
    }
  });
}

