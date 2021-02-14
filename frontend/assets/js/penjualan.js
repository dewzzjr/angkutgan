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
  Form: {
    date: '',
    customer: '',
    project_id: 0,
    deposit: 0,
    address: '',
    shipping_fee: 0,
    items: []
  },
  Set: function (data) {
    Sales.Form = data;
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
  }
};

$(document).ready(function () {
  header();

  var rowIdx = 0;
  var totalPrice = 0;
  var ongkir =0;
  let itemObj = [];
  let projs = [];

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

  function dateFormat(date){
    var dateAr = date.split('-');
    var newDate = dateAr[2] + '/' + dateAr[1] + '/' + dateAr[0];
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

  function updateOngkir($param){
    var newOngkir = double($param.text());
    if (ongkir < newOngkir) {
      totalPrice += (newOngkir - ongkir);
    } else {
      totalPrice -= ongkir - newOngkir;
    }
    ongkir = newOngkir;
  }

  function updateTotal($param){    
    totalPrice += double($param);
  }

  // Display list of item in Ringkasan
  function showRingkasan(){
    totalPrice = ongkir;
    var newList = [];
    $.each(itemObj, function(i, item) {
      var li = (`
      <li class="list-group-item text-right" id="${item["id"]}">
      ${item["code"]} &times; ${item["count"]} ${item["unit"]} = ${item["total"]}
      </li>`);
      newList.push(li);
      updateTotal(item["total"]);
    });
    $('#ringkasanItem span').html(newList.join(''));
    console.log(newList);
    $('#total span').text(currency(totalPrice));
  }
  
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
  $('.autocomplete').on('autocomplete.select', (e, item) => {
    Pelanggan.GetDetail(item.value, (p) => {
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
    if (!$('#nameEdit').val() || !$('#countBarang').val()) {
      $('#tambah').addClass('btn-warning');
      $('#tambah').removeClass('btn-success');
      $('#tambah').text('Gagal');
      setTimeout (function(){
        $('#tambah').addClass('btn-success');
        $('#tambah').removeClass('btn-warning');
        $('#tambah').text('Tambah');
      }, 2000);

      e.preventDefault();
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

      console.log(itemObj);
      console.log(totalPrice);

      $('#addCodeBarang').val('').focus();
      $('#nameEdit').val('');
      $('#countBarang').val('');
      $('#unit span').text('');
      $('#priceBarang').val('');
      $('#totalBarang').attr('value', '');

      e.preventDefault();
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
        $('#editTotal').attr('value', currency((double($(this).val())*double(itemObj[index].total))));
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

      e.preventDefault();
    }
  });

  // Update ringkasan ongkir
  $("#deliveryFee").on("change paste", function() {
    var ongkir = double($(this).val());
    if (isNaN(ongkir)) {
      $('#ongkir span').text(currency(0));
    } else {
      $('#ongkir span').text(currency(ongkir));
    }
    updateOngkir($('#ongkir span'));
    $('#total span').text(currency(totalPrice));
  });

  // Submit transaksi
  $('#submit').on('click', function (e) {
    if (!$('#datePicker').val() || !$('#customerName').val() || itemObj.length < 1) {
      $('#warningSubmit').show('fade');
      setTimeout (function(){
        $('#warningSubmit').hide('fade');
      }, 3000);
    } else {
      var date = $('#datePicker').val();
      var customer = $('#customerName').val();
      var address = $('#customerAddress').val();
      var shipping = ongkir;
      var totalprice = totalPrice;
      var item = [];

      for (var i = 0; i < itemObj.length; i++){
        let listItem = {};
        var id = itemObj[i].id;
        var code = itemObj[i].code;
        var name = itemObj[i].name;
        var price = double(itemObj[i].price);
        var amount = parseInt(itemObj[i].count);
        //var unit = itemObj[i].unit;
        //var total = double(itemObj[i].total);

        listItem ["id"] = id;
        listItem ["code"] = code;
        listItem ["name"] = name;
        listItem ["price"] = price;
        listItem ["amount"] = amount;
        //listItem ["unit"] = unit;
        //listItem ["total"] = total;

        item.push(listItem);
      }

      var items = item;

      Sales.Form ["date"] = dateFormat(date);
      Sales.Form ["customer"] = customer;
      Sales.Form ["address"] = address;
      Sales.Form ["shipping_fee"] = shipping;
      Sales.Form ["totalprice"] = totalprice;
      Sales.Form ["items"] = items;

      console.log(Sales.Form);

      Sales.Create(() => {
        $('#successSubmit').show('fade');
        setTimeout (function(){
          $('#successSubmit').hide('fade');
        }, 3000);
        Loading.End();

        $('#datePicker').empty();
        $('#customerCode').val('');
        $('#customerName').val('');
        $('#formLokasi').empty();
        $('#deliveryFee').val('');

        itemObj = [];
        addTable();
        showRingkasan();

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

});