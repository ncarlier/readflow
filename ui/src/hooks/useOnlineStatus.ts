import { useEffect, useState } from 'react'

import { getOnlineStatus } from '../helpers'

export const useOnlineStatus = () => {
  const [onlineStatus, setOnlineStatus] = useState(getOnlineStatus())

  const goOnline = () => setOnlineStatus(true)

  const goOffline = () => setOnlineStatus(false)

  useEffect(() => {
    window.addEventListener('online', goOnline)
    window.addEventListener('offline', goOffline)

    return () => {
      window.removeEventListener('online', goOnline)
      window.removeEventListener('offline', goOffline)
    }
  }, [])

  return onlineStatus
}
