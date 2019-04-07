import { useState, useEffect, SetStateAction, RefObject } from 'react'

export default (ref: RefObject<HTMLElement>, fetchMoreItems: () => void): [boolean, (value: SetStateAction<boolean>) => void] => {
  const [isFetching, setIsFetching] = useState(false)

  useEffect(() => {
    if (!ref.current) {
      return
    }
    const $container = ref.current.parentElement
    if ($container) {
      $container.addEventListener('scroll', handleScroll)
      $container.addEventListener('resize', handleScroll)
      return () => {
        $container.removeEventListener('scroll', handleScroll)
        $container.removeEventListener('resize', handleScroll)
      }
    }
  }, [])

  useEffect(() => {
    if (!isFetching) return
    fetchMoreItems()
  }, [isFetching])

  function handleScroll() {
    if (isFetching || !ref.current || !ref.current.parentElement) {
      return
    }
    const $container = ref.current.parentElement
    const delta = $container.scrollHeight - $container.scrollTop - $container.offsetHeight
    if (delta < 42) {
      setIsFetching(true)
    }
  }

  return [isFetching, setIsFetching]
}
