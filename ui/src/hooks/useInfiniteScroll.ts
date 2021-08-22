import { RefObject, useCallback, useEffect, useState } from 'react'

export const useInfiniteScroll = (
  ref: RefObject<HTMLElement>,
  hasMore: boolean,
  fetchMoreItems: () => Promise<void>
): boolean => {
  const [isFetching, setIsFetching] = useState(false)

  const fetchMore = useCallback(async () => {
    setIsFetching(true)
    try {
      await fetchMoreItems()
    } finally {
      setIsFetching(false)
    }
  }, [fetchMoreItems])

  const handleScroll = useCallback(() => {
    if (isFetching || !ref.current || !ref.current.parentElement) {
      return
    }
    const $container = ref.current.parentElement
    const delta = $container.scrollHeight - $container.scrollTop - $container.offsetHeight
    if (!isFetching && delta < 42) {
      fetchMore()
    }
  }, [isFetching, fetchMore, ref])

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
  }, [ref, hasMore, handleScroll])

  return isFetching
}
