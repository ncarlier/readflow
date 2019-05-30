import React, { RefObject, useCallback, useEffect, useRef, useState } from 'react'

import ButtonIcon from '../../common/ButtonIcon'
import Center from '../../common/Center'
import Empty from '../../common/Empty'
import Panel from '../../common/Panel'
import { useMedia } from '../../hooks'
import useInfiniteScroll from '../../hooks/useInfiniteScroll'
import useKeyboard from '../../hooks/useKeyboard'
import { Article } from '../models'
import styles from './ArticleList.module.css'
import SwipeableArticleCard from './SwipeableArticleCard'
import ArticleCard from './ArticleCard'

interface Props {
  articles: Article[]
  emptyMessage: string
  hasMore: boolean
  fetchMoreArticles: () => Promise<void>
  refetch: () => Promise<any>
  filter?: (a: Article) => boolean
}

const useKeyNavigation = (ref: RefObject<HTMLUListElement>, itemClassName: string, enable = true) => {
  const left = ['left', 'ArrowLeft', 'p', 'j']
  const right = ['right', 'ArrowRight', 'n', 'k']
  useKeyboard(
    [...left, ...right],
    e => {
      if (ref.current) {
        const { activeElement } = document
        if (activeElement) {
          if (activeElement.className == itemClassName) {
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
  const { fetchMoreArticles, refetch, hasMore, filter = () => true, emptyMessage = 'No more article to read' } = props

  const ref = useRef<HTMLUListElement>(null)
  const [loading, setLoading] = useState(false)
  const [articles, setArticles] = useState(props.articles.filter(filter))
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

  useEffect(() => {
    setArticles(props.articles.filter(filter))
  }, [props.articles])

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
          {isMobileDisplay && !article.isOffline ? (
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
