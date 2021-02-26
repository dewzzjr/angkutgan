const Pelanggan = {
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

const Barang = {
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

const Sales = {
  Form: {},
  Set: function (data) {
    Sales.Form = data;
  },
  Edit: function (callback, failedCallback) {
    let data = this.Form;
    let set = this.Set;
    let retries = false;
    $.ajax({
      type: 'PATCH',
      url: `/sales`,
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
      url: `/sales`,
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
  Create: function (callback, failedCallback) {
    let data = this.Form;
    let set = this.Set;
    $.ajax({
      type: 'POST',
      url: '/sales',
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
  GetDetail: function (customer, date, callback) {
    $.ajax({
      type: 'GET',
      url: `/sales/${customer}/${date}`,
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

  // BUAT TRANSAKSI
  var rowIdx = 0;
  var rowIdxEdit = 0;
  var totalPrice = 0;
  var ongkir = 0;
  var disc = 0;
  var finalPrice = 0;
  let itemObj = [];
  let projs = [];
  var deposit = 0;
  var project_id = 0;

  // Convert number to currency
  function currency($param){
    var number = new Intl.NumberFormat('id-ID', { 
        style: 'currency', 
        currency: 'IDR' 
    }).format($param);
    return number;
  }

  // Convert currency to integer
  function double(currency){  
    var temp = currency.replace(/[^0-9,-]+/g,""); 
    return parseInt(temp); 
  }

  // Convert date format
  function dateFormat(date){
    var dateAr = date.split('-');
    var newDate = dateAr[2] + '/' + dateAr[1] + '/' + dateAr[0];
    return newDate;
  }

  function dateFormatReset(date){
    var dateAr = date.split('/');
    var newDate = dateAr[2] + '-' + dateAr[1] + '-' + dateAr[0];
    return newDate;
  }

  // Add Barang to table
  function addTable(){
    var last = itemObj.slice(-1)[0];
    var tr = $('<tr>');
    var th = $('<th>');
    var tdJumlah = $('<td>');
    var tdHarga = $('<td>');
    var tdTotal = $('<td>');
    var tdBtn = $('<td>');
    th.append(last.code);
    tdJumlah.append(last.count);
    tdHarga.append(last.price);
    tdTotal.append(last.total);
    tdBtn.append(
      `<button type="button" class="btn btn-warning edit">Ubah</button>
      <button type="button" class="btn btn-danger remove">Hapus</button>`
    );
    tr.attr('id', 'item' + (rowIdx));
    tr.append(th);
    tr.append(tdJumlah);
    tr.append(tdHarga);
    tr.append(tdTotal);
    tr.append(tdBtn);
    $('#listItem').append(tr);
  }

  function addTableEdit(){
    var last = itemObj.slice(-1)[0];
    var tr = $('<tr>');
    var th = $('<th>');
    var tdJumlah = $('<td>');
    var tdHarga = $('<td>');
    var tdTotal = $('<td>');
    var tdBtn = $('<td>');
    th.append(last.code);
    tdJumlah.append(last.count);
    tdHarga.append(last.price);
    tdTotal.append(currency(last.count*double(last.price)));
    tdBtn.append(
      `<button type="button" class="btn btn-warning edit">Ubah</button>
      <button type="button" class="btn btn-danger remove">Hapus</button>`
    );
    tr.attr('id', 'item' + (rowIdx));
    tr.append(th);
    tr.append(tdJumlah);
    tr.append(tdHarga);
    tr.append(tdTotal);
    tr.append(tdBtn);
    $('#listItem').append(tr);
  }

  // Update ongkir
  function updateOngkir($param){
    if (ongkir < $param) {
      finalPrice += ($param - ongkir);
    } else {
      finalPrice -= ongkir - $param;
    }
    ongkir = $param;
  }

  // Update total Ringkasan
  function updateTotal($param){    
    totalPrice += double($param);
  }

  function updateTotalEdit($param){    
    totalPrice += double($param);
  }

  function updateFinal(){
    if ((disc == 0) && (ongkir == 0)) {
      finalPrice = totalPrice;
    } else {
      finalPrice = totalPrice + ongkir - disc;
    }
    $('#totalFinal span').text(currency(finalPrice));
  }

  function updateDiscount() {
    var val = $("#discount").val();
    disc = totalPrice * val / 100;
    $('#diskon span').text(`${$("#discount").val()} % = ${currency(disc)}`);
  }

  // Reset form after submit
  function resetForm(){
    $('#datePicker').val('');
    $('#customerCode').val('');
    $('#customerName').val('');
    $('#formLokasi').empty();
    $('#deliveryFee').val(0);
    $('#discount').val(0);
    $('#listItem').empty();
    $('#ringkasan').empty();
    $('#ringkasanItem').empty();
    $('#ongkir span').text('Rp 0,00');
    $('#total span').text('Rp 0,00');
    $('#diskon span').text('0 % = Rp 0,00');
    $('#totalFinal span').text('Rp 0,00');
    $('#submit').removeClass('submitEdit');
    $('#submit').addClass('submit');

    itemObj = [];
  }

  // Display list of item in Ringkasan
  function showRingkasan(){
    var newList = [];
    totalPrice = 0;
    $.each(itemObj, function(i, item) {
      var li = (`
      <li class="list-group-item text-right" id="${item["id"]}">
      ${item["name"]} &times; ${item["count"]} ${item["unit"]} = ${item["total"]}
      </li>`);
      newList.push(li);
      updateTotal(item["total"]);
    });
    $('#ringkasanItem span').html(newList.join(''));
    console.log(newList);
    updateFinal();
    $('#total span').text(currency(totalPrice));
  }

  function showRingkasanEdit(){
    var newList = [];
    totalPrice = 0;
    $.each(itemObj, function(i, item) {
      var li = (`
      <li class="list-group-item text-right" id="${item["id"]}">
      ${item["name"]} &times; ${item["count"]} = ${item["total"]}
      </li>`);
      newList.push(li);
      updateTotalEdit(item["total"]);
    });
    $('#ringkasanItem span').html(newList.join(''));
    console.log(newList);
    updateFinal();
    $('#total span').text(currency(totalPrice));
  }

  // Input diskon
  $("#discount").on('change', function() {
    var max = parseInt($(this).attr('max'));
    var min = parseInt($(this).attr('min'));
    if ($(this).val() > max) {
      $(this).val(max);
    }
    else if ($(this).val() < min) {
      $(this).val(min);
    }
    
    if (isNaN($(this).val()) || !$(this).val()) {
      $('#diskon span').text(`0 % = ${currency(0)}`);
    } else {
      updateDiscount();
      updateFinal();
    }
  });
  
  // Convert number with adding comma per 3 digits
  $('input.number').keyup(function(event) {
    // skip for arrow keys
    if(event.which >= 37 && event.which <= 40) return;
  
    // format number
    $(this).val(function(index, value) {
      return value
      .replace(/\D/g, "")
      .replace(/\B(?=(\d{3})+(?!\d))/g, ".")
      ;
    });
  });

  // AutoComplete Pelanggan
  $('#customerCode').autoComplete({
    resolverSettings: {
      url: '/ajax?action=customers',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });

  // AutoComplete Barang
  $('#addCodeBarang').autoComplete({
    resolverSettings: {
      url: '/ajax?action=items',
      fail: () => {}
    },
    preventEnter: true,
    noResultsText: 'Tidak ditemukan'
  });

  // Selected autocomplete
  $('.autocomplete#customerCode').on('autocomplete.select', (e, customer) => {
    Pelanggan.GetDetail(customer.value, (p) => {
      $('#customerCode').val(p.code);
      $('#customerName').val(p.name);
      
      // Get project
      if (!p.project){
        var address = `
        <label for="customerAddress">Lokasi</label>
        <textarea class="form-control" type="text" id="customerAddress"></textarea>`;
        $('#formLokasi').html(address);
        $('#customerAddress').val(p.address);
      } else {
        var select = `
        <label for="customerAddress">Lokasi</label>
        <select class="form-select form-control col-3 selectAddress">
        <option value="alamat" selected>Alamat</option>`;
        $('#formLokasi').html(select);
        $.each(p.project, function(i, project) {
          $('.selectAddress').append(`<option value="${project["id"]}">${project["name"]}</option>`);
          projs.push(project);
        });
        $('#formLokasi').append(`<textarea class="form-control" type="text" id="customerAddress"></textarea>`);
        $('#customerAddress').val(p.address);
        console.log(projs);

        $(".selectAddress").change(function(){
          for (var i = 0; i < projs.length; i++){
            var value = projs[i].id;
            if ($(this).val() == value) {
              $('#customerAddress').val(projs[i].location);
              $("#customerAddress").attr("disabled", true);
              project_id = value;
            }
          }
          if ($(this).val() == "alamat") {
            $('#customerAddress').val(p.address)
            $("#customerAddress").attr("disabled", false);
          }
        });
      }

      // Get type
      if (p.type == '1') {
        var sum = `
        <h5 class="card-title">${p.name}</h5>
        <h6 class="card-subtitle mb-2 text-muted">${p.address}</h6>
        <p class="card-text">${p.phone}</p>`
      } else {
        var sum = `
        <h5 class="card-title">${p.group_name}</h5>
        <h6 class="card-subtitle mb-2 text-muted">${p.address}</h6>
        <p class="card-text">PIC: ${p.name} - ${p.phone}</p>`
      }
      $('#ringkasan').html(sum);
    });
  });

  $('.autocomplete#addCodeBarang').on('autocomplete.select', (e, item) => {
    Barang.GetDetail(item.value, (b) => {
      $('#setCode').val(b.code);
      $('#nameEdit').val(b.name);
      $('#priceBarang').val(currency(b.price.sell));
      $('#unit span').text(b.unit);
      $("#countBarang").on("change paste keyup", function() {
        if (!$(this).val()) {
          $('#totalBarang').attr('value', currency(0));
        } else {
          $('#totalBarang').attr('value', currency((double($(this).val())*b.price.sell)));
        } 
      });
    });
  });

  // Tambah Barang to list Item
  $('#tambah').on('click', function (e) {     
    e.preventDefault();
    if (!$('#nameEdit').val() || !$('#countBarang').val()) {
      $('#tambah').addClass('btn-warning');
      $('#tambah').removeClass('btn-success');
      $('#tambah').text('Gagal');
      setTimeout (function(){
        $('#tambah').addClass('btn-success');
        $('#tambah').removeClass('btn-warning');
        $('#tambah').text('Tambah');
      }, 2000);
    } else {
      let listItem = {};
      var id = 'item' + (++rowIdx);
      var code = $('#addCodeBarang').val();
      var name = $('#nameEdit').val();
      var price = $('#priceBarang').val();
      var count = $('#countBarang').val();
      var unit = $('#unit span').text();
      var total = $('#totalBarang').val();

      listItem ["id"] = id;
      listItem ["code"] = code;
      listItem ["name"] = name;
      listItem ["price"] = price;
      listItem ["count"] = count;
      listItem ["unit"] = unit;
      listItem ["total"] = total;

      itemObj.push(listItem);
      addTable();
      showRingkasan();
      updateDiscount();
      updateFinal();

      console.log(itemObj);
      console.log(totalPrice);

      $('#addCodeBarang').val('').focus();
      $('#nameEdit').val('');
      $('#countBarang').val('');
      $('#unit span').text('');
      $('#priceBarang').val('');
      $('#totalBarang').attr('value', '');
    }
	});

  // Delete Barang from list Item
	$('#listItem').on('click', '.remove', function () { 
		var tr = $(this).closest('tr');
    tr.remove();
    var id = tr.attr('id');

    var index = itemObj.map(function (item) { return item.id; }).indexOf(id);
    itemObj.splice(index, 1);
    showRingkasan();
    updateDiscount();
    updateFinal();
    console.log(itemObj);
	});

  // Edit Barang from list Item
  $("#listItem").on("click", ".edit", function() {
    var row = $(this).closest('tr');
    var id = row.attr('id');
    var index = itemObj.map(function (item) { return item.id; }).indexOf(id);
    var edit = `
    <form id="formEditItem">
      <table class="table" id="form">
        <tbody>
          <tr>
            <td colspan="6">
              <div class="form-row">
                <div class="form-group col-12 col-md-6">
                  <label for="addCodeBarang">Kode Barang</label>
                  <div class="input-group">
                  <input disabled name="code" class="form-control" id="editCode" value="${itemObj[index].code}">
                  </div>
                </div>
                <div class="form-group col-12 col-md-6">
                  <label for="nameEdit">Nama Barang</label>
                  <input disabled name="name" type="text" class="form-control" id="editName" value="${itemObj[index].name}">
                </div>
                <div class="form-group col-6 col-md-4 align-middle">
                  <label for="priceBarang">Harga</label>
                  <input disabled name="price" type="text" class="form-control" id="editPrice" value="${itemObj[index].price}">
                </div>
                <div class="form-group col-6 col-md-4">
                  <label for="countBarang">Jumlah</label>
                  <div class="input-group">
                    <input name="count" class="form-control number" id="editCount" value="${itemObj[index].count}">
                    <div class="input-group-append" id="editUnit">
                      <span class="input-group-text">${itemObj[index].unit}</span>
                    </div>
                  </div>
                </div>
                <div class="form-group col-6 offset-6 col-md-4 offset-md-0">
                  <label for="totalBarang">Total</label>
                  <div class="input-group">
                    <input disabled name="total" type="text" class="form-control" id="editTotal" value="${itemObj[index].total}">
                    <div class="input-group-append">
                      <button class="btn btn-success update" id="update">Update</button>
                    </div>
                  </div>
                </div>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </form>`;
    row.html(edit);

    $("#editCount").on("change paste keyup", function() {
      if (!$(this).val()) {
        $('#editTotal').attr('value', currency(0));
      } else {
        $('#editTotal').attr('value', currency((double($(this).val())*double(itemObj[index].price))));
      } 
    });
  });

  // Update Barang from list Item
  $("#listItem").on("click", ".update", function(e) {
    if (!$('#editCount').val()) {
      $('#update').addClass('btn-warning');
      $('#update').removeClass('btn-success');
      $('#update').text('Gagal');
      setTimeout (function(){
        $('#update').addClass('btn-success');
        $('#update').removeClass('btn-warning');
        $('#update').text('Tambah');
      }, 2000);

      e.preventDefault();
    } else {
      var row = $(this).closest('tr');
      var id = row.attr('id');
      var index = itemObj.map(function (item) { return item.id; }).indexOf(id);

      itemObj[index].count = $('#editCount').val();
      itemObj[index].total = $('#editTotal').val();

      var update = `
      <th>${itemObj[index].code}</th>
      <td>${itemObj[index].count}</td>
      <td>${itemObj[index].price}</td>
      <td>${itemObj[index].total}</td>
      <td><button type="button" class="btn btn-warning edit">Ubah</button>
        <button type="button" class="btn btn-danger remove">Hapus</button></td>`;
      row.html(update);
      
      showRingkasan();
      updateDiscount();
      updateFinal();

      e.preventDefault();
    }
  });

  // Update ringkasan ongkir
  $("#deliveryFee").on("change paste", function() {
    ongkir = double($(this).val());
    if (isNaN(ongkir)) {
      $('#ongkir span').text(currency(0));
    } else {
      $('#ongkir span').text(currency(ongkir));
      updateOngkir(ongkir);
      updateFinal();
    }
  });

  // Submit transaksi
  $('.submit').off().on('click', function (e) {
    if (!$('#datePicker').val() || !$('#customerName').val() || itemObj.length < 1) {
      $('#warningSubmit').show('fade');
      setTimeout (function(){
        $('#warningSubmit').hide('fade');
      }, 3000);
    } else {
      var date = $('#datePicker').val();
      var customer = $('#customerCode').val();
      var address = $('#customerAddress').val();
      var discount = $('#discount').val();
      var shipping = ongkir;
      var item = [];

      for (var i = 0; i < itemObj.length; i++){
        let listItem = {};
        var code = itemObj[i].code.split('-');
        var newCode = code[0];
        var price = double(itemObj[i].price);
        var amount = parseInt(itemObj[i].count);

        listItem ["code"] = newCode;
        listItem ["price"] = price;
        listItem ["amount"] = amount;

        item.push(listItem);
      }

      var items = item;

      Sales.Form ["date"] = dateFormat(date);
      Sales.Form ["customer"] = customer;
      Sales.Form ["address"] = address;
      Sales.Form ["shipping_fee"] = shipping;
      Sales.Form ["items"] = items;
      Sales.Form ["discount"] = parseInt(discount);
      Sales.Form ["project_id"] = project_id;

      console.log(Sales.Form);

      Sales.Create(() => {
        $('#successSubmit').show('fade');
        setTimeout (function(){
          $('#successSubmit').hide('fade');
        }, 3000);
        Loading.End();

        resetForm();
        $('body').scrollTop(0);

      }, () => {
        $('#failedSubmit').show('fade');
        setTimeout (function(){
          $('#failedSubmit').hide('fade');
        }, 3000);
        Loading.End();
      });

      e.preventDefault();
    }
  });
  // END OF BUAT TRANSAKSI

  // DAFTAR
  Daftar.Init('/sales');
  Daftar.GetData(function (data) {
    $('#tablePenjualan tbody').html('');
    data.forEach(s => {
      if (s.customer.type == "1") {
        var name = s.customer.name;
      } else {
        var name = s.customer.group_name;
      }
      var dateAr = s.tx_date.split('/');
      var date = dateAr[2] + dateAr[1] + dateAr[0];
      var customer = s.customer.code;

      if (s.status.desc == "BARU") {
        var status = `<td>Belum Dibayar</td>`
        var button = `
        <td>
            <button type="button" class="btn btn-info print">Cetak Kwitansi Pembayaran
                Pelanggan</button>
            <br>
            <button type="button" class="btn btn-success paid">Dibayar</button>
            <button type="button" class="btn btn-warning editSales" data-id="DEWANGGA">Ubah</button>
            <button type="button" class="btn btn-danger delete">Hapus</button>
        </td>`
      } else if (s.status.payment_done) {
        var status = `<td>Sudah Dibayar</td>`
        var button = `
        <button type="button" class="btn btn-info print">Cetak Surat Jalan Keluar</button>
        <br>
        <button type="button" class="btn btn-warning shipment">Dikirim</button>`
      } else if (s.status.in_shipping) {
        var status = `<td>Siap Dikirim tgl 20/12/2020</td>`
        var button = `
        <button type="button" class="btn btn-info print">Cetak Surat Jalan Keluar</button>
        <br>
        <button type="button" class="btn btn-warning editShipment">Ubah Waktu Pengiriman</button>`
      } else if (s.status.done) {
        var status = `<td>Selesai</td>`
        var button = `<td></td>`
      }

      var row = `
      <tr data-date=${date} data-customer=${customer}>
        <td>${s.tx_date}</td>
        <td>${name}</td>
        ${status}
        ${button}
      </tr>`
      $('#tablePenjualan').append(row);
    });
  });
  
  function printAlertBox(){
    $('#printAlertBox').show('fade');
    setTimeout (function(){
      $('#printAlertBox').hide('fade');
    }, 3000);
  }
  function paidAlertBox(){
    $('#paidAlertBox').show('fade');
    setTimeout (function(){
      $('#paidAlertBox').hide('fade');
    }, 3000);
  }
  function delAlertBox(){
    $('#delAlertBox').show('fade');
    setTimeout (function(){
      $('#delAlertBox').hide('fade');
    }, 3000);
  }

  $('#tablePenjualan').delegate('.paid', 'click', function (e) {
    paidAlertBox();    
  });

  $('#tablePenjualan').delegate('.print', 'click', function (e) {
    printAlertBox();    
  });

  $('#tablePenjualan').delegate('.delete', 'click', function (e) {
    delAlertBox();    
  });

  // Edit Transaksi
  $('#tablePenjualan').delegate('.editSales', 'click', function (e) {
    var row = $(this).closest('tr');
    var customer = row.data('customer');
    var date = row.data('date');
    var query = {
      customer: customer,
      date: date,
      action: 'create'
    };
    var url = window.location.pathname + '?' + $.param(query);
    console.log(query);
    console.log(url);
    window.location.replace(url);
  });

  var customer = Menu.Query['customer'];
  var date = Menu.Query['date'];
  if (customer && date) {
    $('#submit').removeClass('submit');
    $('#submit').addClass('submitEdit');
    Sales.GetDetail(customer, date, (s) => {
      $('#datePicker').val(dateFormatReset(s.tx_date));
      $('#customerCode').val(s.customer.code);
      $('#customerName').val(s.customer.name);
      disc = s.discount;
      $('#discount').val(disc);
      ongkir = s.shipping_fee;
      $('#deliveryFee').val(ongkir);
      deposit = s.deposit;

      if ($('#customerCode').val()) {
        Pelanggan.GetDetail(s.customer.code, (p) => {
          // Set address
          if (!p.project){
            var address = `
            <label for="customerAddress">Lokasi</label>
            <textarea class="form-control" type="text" id="customerAddress"></textarea>`;
            $('#formLokasi').html(address);
            $('#customerAddress').val(s.address);
          } else {
            var select = `
            <label for="customerAddress">Lokasi</label>
            <select class="form-select form-control col-3 selectAddress">
            <option value="alamat" selected>Alamat</option>`;
            $('#formLokasi').html(select);
            $.each(p.project, function(i, project) {
              if (project["location"] == s.address) {
                $('.selectAddress').append(`<option value="${project["id"]}" selected>${project["name"]}</option>`);
                $('#formLokasi').append(`<textarea class="form-control" type="text" id="customerAddress"></textarea>`);
                $('#customerAddress').val(s.address).attr('disabled', true);
              } else {
                $('.selectAddress').append(`<option value="${project["id"]}">${project["name"]}</option>`);
                $('#formLokasi').append(`<textarea class="form-control" type="text" id="customerAddress"></textarea>`);
                $('#customerAddress').val(s.address);
              }
              projs.push(project);
            });
    
            $(".selectAddress").change(function(){
              for (var i = 0; i < projs.length; i++){
                var value = projs[i].id;
                if ($(this).val() == value) {
                  $('#customerAddress').val(projs[i].location);
                  $("#customerAddress").attr("disabled", true);
                }
              }
              if ($(this).val() == "alamat") {
                $('#customerAddress').val(s.address)
                $("#customerAddress").attr("disabled", false);
              }
            });
          }
    
          // Get type
          if (p.type == '1') {
            var sum = `
            <h5 class="card-title">${p.name}</h5>
            <h6 class="card-subtitle mb-2 text-muted">${s.address}</h6>
            <p class="card-text">${p.phone}</p>`
          } else {
            var sum = `
            <h5 class="card-title">${p.group_name}</h5>
            <h6 class="card-subtitle mb-2 text-muted">${s.address}</h6>
            <p class="card-text">PIC: ${p.name} - ${p.phone}</p>`
          }
          $('#ringkasan').html(sum);
        });
      };

      for (var i = 0; i < s.items.length; i++){
        let listItem = {};
        var id = 'item' + rowIdxEdit;
        var code = s.items[i].code;
        var name = s.items[i].name;
        var price = s.items[i].price;
        var count = s.items[i].amount;
        var total = s.items[i].price * s.items[i].amount;

        listItem ["id"] = id;
        listItem ["code"] = code;
        listItem ["name"] = name;
        Barang.GetDetail(code, (b) => {
          listItem ["unit"] = b.unit;
        });
        listItem ["price"] = currency(price);
        listItem ["count"] = count;
        listItem ["total"] = currency(total);

        itemObj.push(listItem);

        addTableEdit();

        ++rowIdxEdit;
      }

      if ($('#discount').val()) { updateDiscount() };
      $('#ongkir span').text(currency(ongkir));
      showRingkasanEdit();

      console.log(projs);
      console.log(itemObj);
      console.log(ongkir);
      console.log(disc);
    });
  }

  $('.submitEdit').off().on('click', function (e) {
    if (!$('#datePicker').val() || !$('#customerName').val() || itemObj.length < 1) {
      $('#warningSubmit').show('fade');
      setTimeout (function(){
        $('#warningSubmit').hide('fade');
      }, 3000);
    } else {
      var date = $('#datePicker').val();
      var customer = $('#customerCode').val();
      var address = $('#customerAddress').val();
      var discount = $('#discount').val();
      var shipping = ongkir;
      var item = [];

      for (var i = 0; i < itemObj.length; i++){
        let listItem = {};
        var code = itemObj[i].code.split('-');
        var newCode = code[0];
        var price = double(itemObj[i].price);
        var amount = parseInt(itemObj[i].count);

        listItem ["code"] = newCode;
        listItem ["price"] = price;
        listItem ["amount"] = amount;

        item.push(listItem);
      }

      var items = item;

      Sales.Form ["date"] = dateFormat(date);
      Sales.Form ["customer"] = customer;
      Sales.Form ["address"] = address;
      Sales.Form ["shipping_fee"] = shipping;
      Sales.Form ["items"] = items;
      Sales.Form ["deposit"] = deposit;
      Sales.Form ["discount"] = parseInt(discount);

      console.log(Sales.Form);

      Sales.Edit(() => {
        $('#successSubmit').show('fade');
        setTimeout (function(){
          $('#successSubmit').hide('fade');
        }, 3000);
        Loading.End();

        resetForm();
        $('body').scrollTop(0);

      }, () => {
        $('#failedSubmit').show('fade');
        setTimeout (function(){
          $('#failedSubmit').hide('fade');
        }, 3000);
        Loading.End();
      });

      e.preventDefault();
    }
  });
  // END OF DAFTAR
});