const Kirim = {
  Form: {},
  Tx: {},
  SetFunc: (func) => {
    Kirim.GetTx = func;
  },
  Create: (callback, failedCallback) => {
    if (!Kirim.Tx.id) {
      return
    }
    let id = Kirim.Tx.id;
    let data = Kirim.Form;
    $.ajax({
      type: 'POST',
      url: `/shipment/${id}`,
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
    for (let i = 0; i < Kirim.Tx.items.length; i++) {
      let e = Kirim.Tx.items[i];
      if (e.code == code) {
        return e;
      }
    }
    return {};
  },
  Reload: () => {
    let format = Kirim.Tx.tx_date.split("/");
    let date = format[2] + format[1] + format[0];
    var query = {
      customer: Kirim.Tx.customer.code,
      date: date,
      action: 'shipment'
    };
    var url = window.location.pathname + '?' + $.param(query);
    window.location.replace(url);
  }
}

Kirim.Init = function () {
  var getTransaction = function () {
    let customer = $('#shipment [name="customer"]').val();
    let date = $('#shipment [name="tx_date"]').val();
    if (customer && date) {
      Kirim.GetTx(customer, date, (ok) => {
        if (!ok.id) {
          $('#addShipment').attr('disabled', 'disabled');
          $('#dateShipment').attr('disabled', 'disabled');
          return
        }
        let name = ok.customer.group_name ? ok.customer.group_name : ok.customer.name;
        $('#customerNameShipment').val(name);
        $('#addressShipment').val(ok.address);
        $('#shipmentModal .modal-title').html(ok.tx_date + " " + name);

        $('#tableKirim tbody').empty();
        if (ok.status.shipping_done) {
          $('#addShipment').prop('disabled', true);
        } else {
          Kirim.Tx = ok;  
          ok.items.forEach((e) => {
            let status;
            if (e.need_shipment == e.amount) {
              status = 'BELUM DIKIRIM';
            } else if (e.need_shipment == 0) {
              status = 'SUDAH DIKIRIM';
            } else {
              status = `TERKIRIM (${e.amount - e.need_shipment})`;
            }

            let input = '';
            if (status != 'SUDAH DIKIRIM') {
              input = `
              <div class="input-group">
                <div class="input-group-prepend">
                  <span class="input-group-text">
                    MAX ${e.need_shipment}
                  </span>
                </div>
                <input type="number" class="form-control shipAmount" 
                  max="${e.need_shipment}" 
                  data-id="${e.id}"
                  data-code="${e.code}" 
                  value="0">
              </div>`
            }
            let row = `
            <tr>
              <th scope="row">${e.code}</th>
              <td>${e.name}</td>
              <td>${e.amount} ${e.item_unit}</td>
              <td>${status}</td>
              <td>
                ${input}
              </td>
            </tr>`
            $('#tableKirim tbody').append(row);
          });
        }

        $('#tableShipment tbody').empty();
        ok.shipment.forEach((e, idx, arr) => {
          let last = (arr.length - 1) == idx;
          let date = e.date;
          let format = date.split("/");
          let dateID = format[2] + format[1] + format[0];
          let id = `SHIP${dateID}`;
          let items = '';
          let counts = '';
          let deleteBtn = last ? `
          <button type="button" class="btn btn-danger deleteShipment" data-date="${date}">Hapus</button>
          ` : '';
          e.items.forEach((i) => {
            let obj = Kirim.GetItem(i.code);
            let item = `
            <p>${i.code} - ${obj.name}</p>`;
            let count = `
            <p>&times;${i.amount} ${obj.item_unit}</p>`;
            items += item;
            counts += count;
          });
          let row = `
          <tr>
            <td scope="row">${date}</td>
            <td>
              <div class="btn-group">
                <button type="button" class="btn btn-success collapsed" 
                  data-toggle="collapse" 
                  data-target="#${id}"
                  data-parent="#tableShipment">Detail</button>
                ${deleteBtn}
                <button type="button" class="btn btn-info">
                  Cetak
                </button>
              </div>
            </td>
          </tr>
          <tr class="collapse" id="${id}">
            <td class="text-right" scope="row" >
              ${items}
            </td>
            <td>
              ${counts}
            </td>
          </tr>`
          $('#tableShipment tbody').append(row);
        });
        $('#tableShipment .deleteShipment').on('click', function () {
          let id = Kirim.Tx.id;
          let data = {
            date: $(this).data('date'),
          }
          $.confirm({
            title: 'Peringatan!',
            content: `Apakah anda yakin untuk menghapus pengiriman terakhir di transaksi ini?`,
            buttons: {
              ok: function () {
                $.ajax({
                  type: 'DELETE',
                  url: `/shipment/${id}`,
                  data: JSON.stringify(data),
                  contentType: 'application/json',
                  success: function (data, status, xhr) {
                    if (status === 'success' && data.result == 'OK') {
                      $.alert({
                        title: 'Berhasil',
                        content: `berhasil menghapus pengiriman`,
                        onClose: function () {
                          Kirim.Reload();
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
        $('#dateShipment').removeAttr('disabled');
        $('#dateShipment').datepicker({
          format: 'dd/mm/yyyy',
          startDate: '05/06/2021',
        }).on('changeDate', function () {
          let ok = Kirim.Tx;
          let date = $('#shipment [name="shipment_date"]').val();
          if (!ok.id) {
            return
          }
          let name = ok.customer.group_name ? ok.customer.group_name : ok.customer.name;
          $('#addShipment').removeAttr('disabled');
          $('#shipmentModal .modal-title').html(date + " " + name);
        }).on('clearDate', function () {
          $('#addShipment').attr('disabled', 'disabled');
        });
      });
    }
  }

  $('#shipment [name="tx_date"]').datepicker({
    format: 'dd/mm/yyyy',
  }).on('change', function () {
    getTransaction();
  });

  $('#customerCodeShipment').autoComplete({
    resolverSettings: {
      minLength: 2,
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  }).on('autocomplete.select', function () {
    getTransaction();
  });

  if (Menu.Query['date'] && Menu.Query['customer'] && Menu.Query["action"] == 'shipment') {
    var date = Menu.Query["date"];
    date = [date.substring(0, 4), date.substring(4, 6), date.substring(6, 8)].join('-');
    $('#shipment [name="tx_date"]').datepicker('update', new Date(date));
    let cust = Menu.Query["customer"];
    $('#customerCodeShipment').autoComplete('set', {
      value: cust,
      text: cust
    });
    getTransaction();
  } else {
    $('#shipment [name="tx_date"]').datepicker('update', new Date());
  }
  $(document).on('show.bs.modal', '#paymentModal', function (e) {
    console.log('works');
  });

  var serializeObject = function () {
    let items = [];
    $('#formShipment').find('.shipAmount').each((i, input) => {
      let amount = parseInt($(input).val()) || 0;
      if (!amount) {
        return;
      }
      let code = $(input).data('code');
      let itemID = $(input).data('id');
      items.push({
        amount: amount,
        item_id: itemID,
        code: code,
      });
    });
    return {
      date: $('#shipment [name="shipment_date"]').val(),
      items: items,
    }
  }
  $('#formShipment').on('submit', (e) => {
    e.preventDefault();
    let data = serializeObject();
    console.log(data);
    Kirim.Form = data;

    Loading.Start($('#formShipment [type="submit"]'));
    Kirim.Create(() => {
      Loading.End();
      Form.Reset($('#formShipment'));
      Kirim.Reload();
    }, () => {
      Loading.End();
    });
  });
}