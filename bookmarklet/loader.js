function l(d, u, c) {
  let done = false
  const js = d.createElement('script')
  js.type = 'text/javascript'
  js.src = u
  js.onerror = () => alert('Sorry, unable to load bookmarklet.')
  js.onload = js.onreadystatechange = function () {
    if (!done && (!this.readyState || this.readyState == 'loaded' || this.readyState == 'complete')) {
      done = true
      c()
    }
  }
  d.body.appendChild(js)
}
l(document, '//readflow.app/bookmarklet.js', function () {
  window.rfB.boot('/', '//api.readflow.app', 'mykey')
})
