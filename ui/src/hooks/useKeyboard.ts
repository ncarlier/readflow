import mousetrap from 'mousetrap'
import { useEffect } from 'react'

export type KeyHandler = (e: ExtendedKeyboardEvent, combo: string) => void

export default (key: string | string[], handler: KeyHandler, enable = true) => {
  useEffect(() => {
    if (enable) {
      // console.log('bind', key)
      mousetrap.unbind(key)
      setTimeout(() => {
        mousetrap.bind(key, handler)
      }, 200)
      return () => {
        // console.log('unbind', key)
        mousetrap.unbind(key)
      }
    }
  }, [key, handler, enable])
}
