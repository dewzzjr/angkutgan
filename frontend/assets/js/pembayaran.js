const Bayar = {
  Form: {},
  Tx: {},
  SetFunc: (func) => {
    Bayar.GetTx = func;
  },
  Create: (callback, failedCallback) => {
    if (!this.Tx.id) {
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
    if (!this.Tx.id) {
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
  }
}

Bayar.Init = function() {
  var getTransaction = function() {
    let customer = $('#create [name="customer"]').val();
    let date = $('#payment [name="tx_date"]').val();
    Bayar.GetTx(customer, date, (ok) => {
      console.log(ok)
    });
  }
  $('#payment [name="tx_date"]').datepicker({
    format: 'dd/mm/yyyy',
  }).on('change', getTransaction);

  $('#payment [name="customer"]').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });
  $('#payment [name="customer"]').on('autocomplete.select', getTransaction);

  $('#payment .addPayment').on('click', () => {
    $('#paymentModal [name="date"]').datepicker({ format: 'dd/mm/yyyy' });
    $('#paymentModal [type="submit"]').data('action', 'add');
  });
  $('#payment .editPayment').on('click', () => {
    $('#paymentModal [name="date"]').datepicker({ format: 'dd/mm/yyyy' });
    $('#paymentModal [type="submit"]').data('action', 'edit');
  });
  $('#formPayment').on('submit', (e) => {
    e.preventDefault();
    let action = $('#paymentModal [type="submit"]').data('action');
    Bayar.Form = {
      date: $('#formPayment [name="date"]').val(),
      name: $('#formPayment [name="name"]').val(),
      method: parseInt($('#formPayment [name="method"]').val()) || 100,
      amount: parseInt($('#formPayment [name="amount"]').val()) || 0,
    }
    Form.Reset('#paymentIndex');
    if (action == 'edit') {
      Bayar.Edit(() => {
        Form.Message('success', 'berhasil mengubah pembayaran', $('#messagePay'));
        Form.Reset($('#formPayment'));
      }, () => {
        Form.Message('danger', 'gagal mengubah pembayaran', $('#messagePay'));
        Loading.End();
      });
    } else if (action == 'add') {
      Bayar.Create(() => {
        Form.Message('success', 'berhasil menambah pembayaran', $('#messagePay'));
        Form.Reset($('#formPayment'));
      }, () => {
        Form.Message('danger', 'gagal menambah pembayaran', $('#messagePay'));
      });
    }
  });
}

