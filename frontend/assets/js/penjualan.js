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
  //header();
  
  function currency($param){
    return 'Rp'+($param.toLocaleString());
  }
  
  $('input.number').keyup(function(event) {
    // skip for arrow keys
    if(event.which >= 37 && event.which <= 40) return;
  
    // format number
    $(this).val(function(index, value) {
      return value
      .replace(/\D/g, "")
      .replace(/\B(?=(\d{3})+(?!\d))/g, ",")
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
        $('#totalBarang').attr('value', currency(($(this).val().replace(",", "")*b.price.sell))); 
      });

      $('#tambah').click(function(){
        $("#tableBarang tbody").prepend(
          `<tr>
            <th scope="row">${b.code}</th>
            <td>${$("#countBarang").val()} ${b.unit}</td>
            <td>${currency(b.price.sell)}</td>
            <th>${$('#totalBarang').val()}</th>
            <td>
                <button type="button" class="btn btn-warning">Ubah</button>
                <button type="button" class="btn btn-danger">Hapus</button>
            </td>
          </tr>`
        );

        $('#setCode').val('');
        $('#nameEdit').val('');
        $('#priceBarang').val('');
        $('#unit span').val('');
        $("#countBarang").val('');
        $('#totalBarang').attr('value', '');
      });
    });
  });

  /* $("#deliveryFee").on("change paste keyup", function() {
    $('#ongkir').text('Ongkos Kirim = Rp' + $(this).val());
  });

  if ($("input:empty").length == 0){
    $('#total').text('Total Tagihan = Rp' + $('#ongkir').val());
  }*/ 


});