function urlBase64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
  const rawData = window.atob(base64)
  return Uint8Array.from(Array.from(rawData).map(char => char.charCodeAt(0)))
}

export const setupNotification = () => {
  if (!('Notification' in window)) {
    console.error('This browser does not support desktop notification')
  } else if (Notification.permission === 'granted') {
    console.log('Permission to receive notifications has been granted')
  } else if (Notification.permission !== 'denied') {
    Notification.requestPermission(function(permission) {
      if (permission === 'granted') {
        console.log('Permission to receive notifications has been granted')
      }
    })
  }
}

const renewPushSubscription = async (registration: ServiceWorkerRegistration) => {
  const applicationServerKey = urlBase64ToUint8Array(window.vapidPublicKey)
  try {
    const subscripiton = await registration.pushManager.getSubscription()
    if (subscripiton) {
      const ok = await subscripiton.unsubscribe()
      if (ok) {
        await registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey
        })
      }
    }
  } catch (e) {
    console.error(e)
  }
}

export const subscribePush = async (registration: ServiceWorkerRegistration) => {
  const applicationServerKey = urlBase64ToUint8Array(window.vapidPublicKey)

  try {
    const subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey
    })
    console.log('Subscribed to push manager:', subscription.endpoint)
  } catch (err) {
    console.error(err)
    // Try to renew subscription if allready subscribed
    // with a different applicationServerKey
    renewPushSubscription(registration)
  }
}
