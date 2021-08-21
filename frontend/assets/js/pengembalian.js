const Kembali = {
  Form: {},
  Tx: {},
  SetFunc: (func) => {
    Kembali.GetTx = func;
  },
  Create: (callback, failedCallback) => {
    if (!Kembali.Tx.id) {
      return
    }
    let id = Kembali.Tx.id;
    let data = Kembali.Form;
    $.ajax({
      type: 'POST',
      url: `/return/${id}`,
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
  GetItem: (code) => {
    for (let i = 0; i < Kembali.Tx.items.length; i++) {
      let e = Kembali.Tx.items[i];
      if (e.code == code) {
        return e;
      }
    }
    return {};
  },
  Reload: () => {
    let format = Kembali.Tx.tx_date.split("/");
    let date = format[2] + format[1] + format[0];
    var query = {
      customer: Kembali.Tx.customer.code,
      date: date,
      action: 'return'
    };
    var url = window.location.pathname + '?' + $.param(query);
    window.location.replace(url);
  }
}

Kembali.Init = function () {
}