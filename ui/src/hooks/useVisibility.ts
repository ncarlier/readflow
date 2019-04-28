import { RefObject, useEffect, useRef, useState } from 'react'

export default (ref: RefObject<HTMLElement>, options?: IntersectionObserverInit) => {
  const [visible, setVisibilty] = useState<IntersectionObserverEntry | null>(null)
  const isIntersecting = useRef(false)

  const handleObserverUpdate: IntersectionObserverCallback = entries => {
    const ent = entries[0]
    if (isIntersecting.current !== ent.isIntersecting) {
      setVisibilty(ent)
      isIntersecting.current = ent.isIntersecting
    }
  }

  const observer = new IntersectionObserver(handleObserverUpdate, options)

  useEffect(() => {
    const $element = ref.current
    if (!$element) {
      return
    }

    observer.observe($element)
    return () => {
      observer.unobserve($element)
    }
  })

  return visible
}
