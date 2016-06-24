$(() => {
  $('.baseurl').text(
    document.location.protocol + '//' + document.location.host
  );

  $('form').on('submit', (event) => {
    event.preventDefault()
    let emails = []
    $('textarea').val().split(/\s/g).forEach((e) => {
      if (e.trim()) {
        emails.push(e.trim())
      }
    })
    var container = $('#results tbody')
    $('tr', container).remove()
    emails.forEach((e) => {
      $('<tr>')
        .append(
          $('<td>').text(e)
        )
        .append(
          $('<td>')
          .append($('<span class="tag">hmm</span>'))
        )
        .appendTo(container)
    })
    $('#results').show()
    if (emails.length) {
      $.ajax({
          url: '/staff',
          data: JSON.stringify({
            email: emails
          }),
          method: 'POST',
          dataType: 'json',
        })
        .then((r) => {
          $('#error').hide()
          $('tr', container).each((i, row) => {
            let email = $('td:first-child', row).text();
            // $('td:last-child', row).text(r[email])
            $('td:last-child span', row).remove()
            if (r[email]) {
              $('td:last-child', row)
              .append($('<span class="tag is-success">Yes!</span>'))
            } else {
              $('td:last-child', row)
              .append($('<span class="tag is-danger">Nope</span>'))
            }
          })
        })
        .fail((err) => {
          // XXX check if the status code was 400
          $('#error pre')
            .text(JSON.stringify(JSON.parse(err.responseText), undefined, 2))
          $('#results').hide()
          $('#error').show()
          console.error(err);
          // alert('Failed')
        })
    }
  })
})
