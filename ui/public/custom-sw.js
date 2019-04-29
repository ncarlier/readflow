/* eslint-disable no-undef */
self.addEventListener('push', function(event) {
  if (!('Notification' in self && Notification.permission === 'granted')) {
    return
  }

  let data = {}
  if (event.data) {
    try {
      data = event.data.json()
    } catch (err) {
      data = { body: event.data.text() }
    }
  }
  let title = data.title || 'Something has happened'
  let body = data.body || "Here's something you might want to check out."
  const icon = 'logo.png'
  const tag = 'readflow-notification'

  event.waitUntil(self.registration.showNotification(title, { body, icon, tag }))
})

self.addEventListener('notificationclick', function(event) {
  console.log('On notification click: ', event.notification.tag)
  event.notification.close()
  event.waitUntil(
    clients
      .matchAll({
        type: 'window'
      })
      .then(function(clientList) {
        for (var i = 0; i < clientList.length; i++) {
          var client = clientList[i]
          if (client.url == '/' && 'focus' in client) return client.focus()
        }
        if (clients.openWindow) {
          return clients.openWindow('/')
        }
      })
  )
})
