import { getAPIURL } from './fetchAPI'

function urlBase64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
  const rawData = window.atob(base64)
  return Uint8Array.from(Array.from(rawData).map((char) => char.charCodeAt(0)))
}

export const isNotificationSupported = () => 'Notification' in window
export const isNotificationGranted = () => isNotificationSupported() && Notification.permission === 'granted'

export const unSubscribePush = async (registration: ServiceWorkerRegistration) => {
  const subscription = await registration.pushManager.getSubscription()
  if (subscription) {
    const ok = await subscription.unsubscribe()
    if (!ok) {
      throw new Error('Unable to renew push manager subscription.')
    }
    return
  }
}

export const subscribePush = async (registration: ServiceWorkerRegistration) => {
  let applicationServerKey: Uint8Array | undefined
  try {
    let subscription = await registration.pushManager.getSubscription()
    if (!subscription) {
      // No subscription: creat a new one
      const res = await fetch(getAPIURL('/info'))
      const data = await res.json()
      applicationServerKey = urlBase64ToUint8Array(data.vapid)
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey,
      })
    }
    console.log('Subscribed to push manager:', subscription.endpoint)
    return subscription
  } catch (err) {
    console.error('Error when creating push subscription:', err)
    throw err
  }
}
