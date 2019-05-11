import mousetrap from 'mousetrap'
import { useEffect } from 'react'

export type KeyHandler = (e: ExtendedKeyboardEvent, combo: string) => void

export default (key: string | string[], handler: KeyHandler, enable = true) => {
  useEffect(() => {
    if (enable) {
      // console.log('bind', key)
      mousetrap.unbind(key)
      mousetrap.bind(key, handler)
      return () => {
        // console.log('unbind', key)
        mousetrap.unbind(key)
      }
    }
  }, [handler, enable])
}
