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
  var getTransaction = function () {
    let customer = $('#return [name="customer"]').val();
    let date = $('#return [name="tx_date"]').val();
    if (customer && date) {
      Kembali.GetTx(customer, date, (ok) => {
        if (!ok.id) {
          $('#addReturn').attr('disabled', 'disabled');
          $('#dateReturn').attr('disabled', 'disabled');
          return
        }
        let name = ok.customer.group_name ? ok.customer.group_name : ok.customer.name;
        $('#customerNameReturn').val(name);
        // $('#addressReturn').val(ok.address);
        $('#returnModal .modal-title').html(ok.tx_date + " " + name);

        $('#tableKembali tbody').empty();
        if (ok.status.shipping_done) {
          $('#addReturn').prop('disabled', true);
        } else {
          Kirim.Tx = ok;
          ok.items.forEach((e) => {
            let status;
            // if (e.need_shipment == e.amount) {
            //   status = 'BELUM DIKIRIM';
            // } else if (e.need_shipment == 0) {
            //   status = 'SUDAH DIKIRIM';
            // } else {
            //   status = `TERKIRIM (${e.amount - e.need_shipment})`;
            // }

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

        $('#tableReturn tbody').empty();
        ok.shipment.forEach((e, idx, arr) => {
          let last = (arr.length - 1) == idx;
          let date = e.date;
          let format = date.split("/");
          let dateID = format[2] + format[1] + format[0];
          let id = `SHIP${dateID}`;
          let items = '';
          let counts = '';
          let deleteBtn = last ? `
          <button type="button" class="btn btn-danger deleteReturn" data-date="${date}">Hapus</button>
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
                  data-parent="#tableReturn">Detail</button>
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
          $('#tableReturn tbody').append(row);
        });
        $('#tableReturn .deleteReturn').on('click', function () {
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
              cancel: function () { }
            },
          });
        });
        $('#dateReturn').removeAttr('disabled');
        $('#dateReturn').datepicker({
          format: 'dd/mm/yyyy',
          startDate: '05/06/2021',
        }).on('changeDate', function () {
          let ok = Kirim.Tx;
          let date = $('#shipment [name="shipment_date"]').val();
          if (!ok.id) {
            return
          }
          let name = ok.customer.group_name ? ok.customer.group_name : ok.customer.name;
          $('#addReturn').removeAttr('disabled');
          $('#shipmentModal .modal-title').html(date + " " + name);
        }).on('clearDate', function () {
          $('#addReturn').attr('disabled', 'disabled');
        });
      });
    }
  }
}