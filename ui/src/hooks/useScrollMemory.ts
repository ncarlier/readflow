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

  const handleScroll = (evt: Event) => {
    if (evt.target) {
      const $el: any = evt.target
      const pos = $el.scrollTop
      const key = location.pathname
      // console.log(`saving scroll position for ${key}: ${pos}`)
      setScrollPosition(key, pos)
    }
  }

  useEffect(() => {
    if (ref.current) {
      ref.current.addEventListener('scroll', handleScroll)
    }
    return () => {
      if (ref.current) {
        ref.current.removeEventListener('scroll', handleScroll)
      }
    }
  })
}
