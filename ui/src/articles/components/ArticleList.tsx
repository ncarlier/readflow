import React, { ReactNode, RefObject, useCallback, useEffect, useRef, useState } from 'react'

import { ButtonIcon, Center, Empty, Panel } from '../../components'
import { useInfiniteScroll, useKeyboard, useMedia } from '../../hooks'
import { Article } from '../models'
import { ArticleCard, SwipeableArticleCard } from '.'
import styles from './ArticleList.module.css'

interface Props {
  articles: Article[]
  empty: ReactNode
  hasMore: boolean
  swipeable?: boolean
  variant?: 'list' | 'grid'
  fetchMoreArticles: () => Promise<void>
  refetch: () => Promise<any>
}

const useKeyNavigation = (ref: RefObject<HTMLUListElement>, itemClassName: string, enable = true) => {
  const left = ['left', 'ArrowLeft', 'p', 'j']
  const right = ['right', 'ArrowRight', 'n', 'k']
  useKeyboard(
    [...left, ...right],
    (e) => {
      if (ref.current) {
        const { activeElement } = document
        if (activeElement) {
          if (activeElement.className === itemClassName) {
            const $el: any = left.includes(e.key) ? activeElement.previousSibling : activeElement.nextSibling
            if ($el) {
              $el.focus()
              return
            }
          }
        }
      }
    },
    enable
  )
}

export const ArticleList = (props: Props) => {
  const {
    articles,
    fetchMoreArticles,
    refetch,
    hasMore,
    empty = <span>No more article to read</span>,
    swipeable = false,
    variant = 'list',
  } = props

  const ref = useRef<HTMLUListElement>(null)
  const [loading, setLoading] = useState(false)
  const [activeIndex, setActiveIndex] = useState(0)

  const isFetchingMore = useInfiniteScroll(ref, hasMore, fetchMoreArticles)
  const isMobileDisplay = useMedia('(max-width: 767px)')

  useKeyNavigation(ref, styles.item, !isMobileDisplay)

  const reload = useCallback(async () => {
    setLoading(true)
    await refetch()
    setLoading(false)
  }, [refetch])

  useEffect(() => {
    if (ref.current) {
      const $el: any = ref.current.childNodes.item(activeIndex)
      if ($el) $el.focus()
    }
  }, [activeIndex])

  useEffect(() => {
    if (!isFetchingMore && !loading && hasMore && articles.length <= 3) {
      setLoading(true)
      fetchMoreArticles().finally(() => setLoading(false))
    }
  }, [isFetchingMore, loading, hasMore, articles, fetchMoreArticles])

  if (articles.length === 0) {
    return (
      <Empty>
        {empty}
        <br />
        <ButtonIcon title="Refresh" icon="refresh" onClick={reload} loading={loading} />
      </Empty>
    )
  }

  return (
    <ul className={styles[variant]} ref={ref}>
      {articles.map((article, idx) => (
        <li key={`article-${article.id}`} className={styles.item} tabIndex={-1} onFocus={() => setActiveIndex(idx)}>
          {swipeable ? (
            <SwipeableArticleCard article={article} />
          ) : (
            <ArticleCard article={article} isActive={idx === activeIndex} />
          )}
        </li>
      ))}
      {(isFetchingMore || loading) && (
        <li>
          <Panel>
            <Center>Fetching more articles ...</Center>
          </Panel>
        </li>
      )}
    </ul>
  )
}
