import { RefObject, useContext, useEffect } from 'react'

import { ScrollMemoryContext, setScrollPosition } from '../context/ScrollMemoryContext'

export default (ref: RefObject<HTMLDivElement>) => {
  const scrollPosition = useContext(ScrollMemoryContext)
  useEffect(() => {
    if (ref.current && scrollPosition > 0) {
      // console.log(`restoring scroll position: ${scrollPosition}`)
      ref.current.scrollTo(0, scrollPosition)
    }
  }, [ref, scrollPosition])

  useEffect(() => {
    const key = window.location.pathname
    const { current } = ref
    return () => {
      if (current) {
        const pos = current.scrollTop
        // console.log(`saving scroll position for ${key}: ${pos}`)
        setScrollPosition(key, pos)
      }
    }
  }, [ref])
}
