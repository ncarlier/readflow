import { RefObject, useContext, useEffect } from 'react'

import { ScrollMemoryContext, setScrollPosition } from '../context/ScrollMemoryContext'

export default (ref: RefObject<HTMLDivElement>) => {
  const scrollPosition = useContext(ScrollMemoryContext)
  useEffect(() => {
    if (ref.current && scrollPosition > 0) {
      // console.log(`restoring scroll position: ${scrollPosition}`)
      ref.current.scrollTo(0, scrollPosition)
    }
  }, [scrollPosition])

  useEffect(() => {
    const key = window.location.pathname
    return () => {
      if (ref.current) {
        const pos = ref.current.scrollTop
        // console.log(`saving scroll position for ${key}: ${pos}`)
        setScrollPosition(key, pos)
      }
    }
  })
}
