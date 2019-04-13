import {useEffect} from "react"
import mousetrap from "mousetrap"

type KeyHandler = (e: ExtendedKeyboardEvent, combo: string) => void

export default (key: string, handler: KeyHandler, enable = true) => {
  useEffect(() => {
    if (enable) {
      mousetrap.unbind(key)
      mousetrap.bind(key, handler)
      return () => {
        mousetrap.unbind(key)
      }
    }
  }, [handler])
}
