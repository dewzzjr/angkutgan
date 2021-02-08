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
  Set: function(data){
    Sales.Form = data;
  },
  Clear: function(callback){
    Sales.Form = {};
    if (callback){
      callback();
    }
  },
  Validate: function (isEdit = false, callback){
    let ok = {
      message: [],
      valid: true
    };

  }
};

$(document).ready(function () {
  header();

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
      $('#customerAddress').val(p.address);
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
        $('#totalBarang').attr('value', currency((double($(this).val())*b.price.sell))); 
      });
    });
  });

  // Tambah Barang to list Item
	var rowIdx = 0;
  $('#tambah').on('click', function () { 
		$('#listItem').append(
	    `<tr id="item${++rowIdx}">
				<th scope="row">${$('#addCodeBarang').val()}</th>
				<td id="listJumlah">${$('#countBarang').val()}</td>
	      <td id="listHarga">${$('#priceBarang').val()}</td>
	      <th id="listTotal">${$('#totalBarang').val()}</th>
	      <td align="right">
	       	<button type="button" class="btn btn-warning edit">Ubah</button>
	       	<button type="button" class="btn btn-danger remove">Hapus</button>
	      </td>
			</tr>`);
    $('#ringkasanItem').append(
      `<li class="list-group-item text-right" id="item${++rowIdx-1}">
      ${$('#nameEdit').val()} &times; ${$('#countBarang').val()} ${$('#unit span').text()} = <span>${$('#totalBarang').val()}</span>
      </li>`);
    $('#addCodeBarang').val('').focus();
    $('#nameEdit').val('');
		$('#countBarang').val('');
    $('#unit span').text('');
	  $('#priceBarang').val('');
	  $('#totalBarang').attr('value', '');
	}); 

  // Delete Barang from list Item
	$('#listItem').on('click', '.remove', function () { 
		var tr = $(this).closest('tr'); 
    var id = tr.attr('id');
    $('[id="'+id+'"]').remove();
		rowIdx--; 
	});

  // Edit Barang from list Item
  $("#listItem").on("click", ".edit", function() {
    var row = $(this).closest('tr');
    var id = row.attr('id');
    if ($('.edit').text() == 'Update'){
      row.find("#listJumlah").text($('#editListJumlah').val());
      row.find("#listHarga").text($('#editListHarga').val());
      row.find("#listTotal").text($('#editListTotal').val());
      $('.edit').text('Edit');
      $('[id="'+id+'"] span').text($('#editListTotal').val());
    } else {
      row.find("#listJumlah").html(`<input class="form-control number" id="editListJumlah">`);
      row.find("#listHarga").html(`<input disabled class="form-control" id="editListHarga" value="${$('#listHarga').text()}">`);
      row.find("#listTotal").html(`<input disabled class="form-control" id="editListTotal" value="${$('#listTotal').text()}">`);
      $('#editListJumlah').focus();
      $('.edit').text('Update');
    }

    $("#editListJumlah").on("change paste keyup", function() {
      var tot = double($(this).val())*double($("#editListHarga").val());
      if (isNaN(tot)) {
        $('#editListTotal').attr('value', currency(0));
      } else {
        $('#editListTotal').attr('value', currency(tot));
      } 
    });
  });

  // Update ringkasan ongkir
  $("#deliveryFee").on("change paste keyup", function() {
    var ongkir = double($(this).val());
    if (isNaN(ongkir)) {
      $('#ongkir').text('Ongkos Kirim = ' + currency(0));
    } else {
      $('#ongkir').text('Ongkos Kirim = ' + currency(ongkir));
    }
  });

});