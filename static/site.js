$(() => {
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
          $('<td>').text('...')
        )
        .appendTo(container)
    })
    $('#results').show();
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
          $('tr', container).each((row) => {
            let email = $('td:first-child', row).text();
            $('td:last-child', row).text(r[email])
          })
        })
        .fail((err) => {
          // XXX check if the status code was 400
          $('#error pre')
            .text(JSON.stringify(JSON.parse(err.responseText), undefined, 2))
          $('#error').show()
          console.error(err);
          // alert('Failed')
        })
    }
  })
})
