import { RefObject, useCallback, useEffect, useState } from 'react'

export default (ref: RefObject<HTMLElement>, hasMore: boolean, fetchMoreItems: () => Promise<void>): boolean => {
  const [isFetching, setIsFetching] = useState(false)

  useEffect(() => {
    if (!isFetching) return
    fetchMoreItems().finally(() => setIsFetching(false))
  }, [isFetching])

  const handleScroll = useCallback(() => {
    if (isFetching || !ref.current || !ref.current.parentElement) {
      return
    }
    const $container = ref.current.parentElement
    const delta = $container.scrollHeight - $container.scrollTop - $container.offsetHeight
    if (delta < 42) {
      setIsFetching(true)
    }
  }, [isFetching, ref])

  useEffect(() => {
    if (!ref.current || !hasMore) {
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
  }, [ref, hasMore])

  return isFetching
}
