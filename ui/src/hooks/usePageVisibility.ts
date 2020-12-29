import { useEffect, useState } from 'react'

const getPageVisibilityApiAttributes = () => {
  let hidden = 'hidden'
  let visibilityChange = 'visibilitychange'
  if ('mozHidden' in document) {
    // Firefox up to v17
    hidden = 'mozHidden'
    visibilityChange = 'mozvisibilitychange'
  } else if ('webkitHidden' in document) {
    // Chrome up to v32, Android up to v4.4, Blackberry up to v10
    hidden = 'webkitHidden'
    visibilityChange = 'webkitvisibilitychange'
  }
  return { hidden, visibilityChange }
}

const isSupported = () => {
  const { hidden } = getPageVisibilityApiAttributes()
  return hidden in document
}

const { hidden, visibilityChange } = getPageVisibilityApiAttributes()

export default () => {
  const [isVisible, setIsVisible] = useState<boolean>(!isSupported() || (document as any)[hidden])

  useEffect(() => {
    if (isSupported()) {
      const handler = () => {
        setIsVisible(!(document as any)[hidden])
      }
      window.addEventListener(visibilityChange, handler)
      return () => {
        window.removeEventListener(visibilityChange, handler)
      }
    }
  }, [])

  return isVisible
}
