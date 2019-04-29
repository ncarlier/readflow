import { API_BASE_URL } from './constants'

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

const renewPushSubscription = async (registration: ServiceWorkerRegistration, applicationServerKey: Uint8Array) => {
  try {
    let subscription = await registration.pushManager.getSubscription()
    if (subscription) {
      const ok = await subscription.unsubscribe()
      if (ok) {
        console.log('Un-subscribed with success. Subscribe again...')
        subscription = await registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey
        })
        console.log('Subscribed to push manager:', subscription.endpoint)
      }
    }
  } catch (err) {
    console.error(err)
  }
}

export const subscribePush = async (registration: ServiceWorkerRegistration) => {
  let applicationServerKey: Uint8Array | undefined
  try {
    let subscription = await registration.pushManager.getSubscription()
    if (!subscription) {
      const res = await fetch(API_BASE_URL)
      const data = await res.json()
      applicationServerKey = urlBase64ToUint8Array(data.vapid)
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey
      })
      console.log('Subscribed to push manager:', subscription.endpoint)
    } else {
      console.log('Push subscription:', subscription.endpoint)
    }
  } catch (err) {
    console.error(err)
    if (applicationServerKey) {
      // Try to renew subscription if allready subscribed
      // with a different applicationServerKey
      console.log('Renewing push subscription...')
      renewPushSubscription(registration, applicationServerKey)
    }
  }
}
