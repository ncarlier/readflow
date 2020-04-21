import React, { RefObject, useCallback, useEffect, useRef, useState } from 'react'

import ButtonIcon from '../../components/ButtonIcon'
import Center from '../../components/Center'
import Empty from '../../components/Empty'
import Panel from '../../components/Panel'
import { useMedia } from '../../hooks'
import useInfiniteScroll from '../../hooks/useInfiniteScroll'
import useKeyboard from '../../hooks/useKeyboard'
import { Article } from '../models'
import ArticleCard from './ArticleCard'
import styles from './ArticleList.module.css'
import SwipeableArticleCard from './SwipeableArticleCard'

interface Props {
  articles: Article[]
  emptyMessage: string
  hasMore: boolean
  swipeable?: boolean
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

export default (props: Props) => {
  const {
    articles,
    fetchMoreArticles,
    refetch,
    hasMore,
    emptyMessage = 'No more article to read',
    swipeable = false,
  } = props

  const ref = useRef<HTMLUListElement>(null)
  const [loading, setLoading] = useState(false)
  const [activeIndex, setActiveIndex] = useState(0)

  const isFetching = useInfiniteScroll(ref, hasMore, fetchMoreArticles)
  const isMobileDisplay = useMedia('(max-width: 767px)')

  useKeyNavigation(ref, styles.item, !isMobileDisplay)

  const refresh = useCallback(async () => {
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

  if (articles.length <= 3) {
    if (hasMore) {
      refetch()
    } else if (articles.length === 0) {
      return (
        <Empty>
          <ButtonIcon title="Refresh" icon="refresh" onClick={() => refresh()} loading={loading} />
          <br />
          <span>{emptyMessage}</span>
        </Empty>
      )
    }
  }

  return (
    <ul className={styles.list} ref={ref}>
      {articles.map((article, idx) => (
        <li key={`article-${article.id}`} className={styles.item} tabIndex={-1} onFocus={() => setActiveIndex(idx)}>
          {swipeable ? (
            <SwipeableArticleCard article={article} />
          ) : (
            <ArticleCard article={article} isActive={idx === activeIndex} />
          )}
        </li>
      ))}
      {isFetching && (
        <li>
          <Panel>
            <Center>Fetching more articles...</Center>
          </Panel>
        </li>
      )}
    </ul>
  )
}
