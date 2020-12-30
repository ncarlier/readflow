import { useEffect, useState } from 'react'

const isSupported = () => {
  return 'visibilityState' in document
}

export default () => {
  const [isVisible, setIsVisible] = useState<boolean>(!isSupported() || document.visibilityState === 'visible')

  useEffect(() => {
    if (isSupported()) {
      const handler = () => {
        setIsVisible(document.visibilityState === 'visible')
      }
      window.addEventListener('visibilitychange', handler)
      return () => {
        window.removeEventListener('visibilitychange', handler)
      }
    }
  }, [])

  return isVisible
}
