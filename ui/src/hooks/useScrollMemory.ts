import { RefObject, useContext, useEffect, useLayoutEffect } from 'react'

import { ScrollMemoryContext, setScrollPosition } from '../contexts/ScrollMemoryContext'

export default (ref: RefObject<HTMLDivElement>) => {
  const scrollPosition = useContext(ScrollMemoryContext)
  useEffect(() => {
    if (ref.current && scrollPosition > 0) {
      // console.log(`restoring scroll position: ${scrollPosition}`)
      ref.current.scrollTo(0, scrollPosition)
    }
  }, [ref, scrollPosition])

  useLayoutEffect(() => {
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
