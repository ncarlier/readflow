/* eslint-disable no-undef */
function b() {
  var r = new Request(FP_URL + '/articles')
  var h = new Headers({
    Accept: 'application/json',
    'Content-Type': 'application/json'
  })
  h.set('Authorization', 'Basic ' + FP_CRED)
  fetch(r, {
    method: 'POST',
    headers: h,
    mode: 'cors',
    body: JSON.stringify({ title: document.title, url: document.location.href })
  }).then(
    function(res) {
      if (res.ok) {
        alert('Webpage added to your readflow.')
      } else {
        alert('Unable to send the web page to readflow: ',  res.statusText)
      }
    },
    function(err) {
      alert('Unable to send the web page to readflow: ', err.message)
    }
  )
}
b()
