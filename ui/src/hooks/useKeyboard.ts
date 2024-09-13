import mousetrap from 'mousetrap'
import { useEffect } from 'react'

export type KeyHandler = (e: mousetrap.ExtendedKeyboardEvent, combo: string) => void

export const useKeyboard = (key: undefined | string | string[], handler: KeyHandler, enable = true) => {
  useEffect(() => {
    if (key && enable) {
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
