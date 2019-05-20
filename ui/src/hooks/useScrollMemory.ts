import { useContext, useEffect } from 'react'

import { ScrollMemoryContext, setScrollPosition } from '../context/ScrollMemoryContext'

export default (target: string) => {
  const scrollPosition = useContext(ScrollMemoryContext)
  useEffect(() => {
    if (scrollPosition > 0) {
      const $el = document.getElementById(target)
      if ($el) {
        $el.scrollTo(0, scrollPosition)
      }
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
    const $el = target ? document.getElementById(target) || window : window
    $el.addEventListener('scroll', handleScroll)
    return () => {
      $el.removeEventListener('scroll', handleScroll)
    }
  })
}
