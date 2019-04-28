import React, { createRef } from 'react'

import ButtonIcon from '../../common/ButtonIcon'
import Center from '../../common/Center'
import Empty from '../../common/Empty'
import Panel from '../../common/Panel'
import { useMedia } from '../../hooks'
import useInfiniteScroll from '../../hooks/useInfiniteScroll'
import { Article } from '../models'
import ArticleCard from './ArticleCard'
import styles from './ArticleList.module.css'
import SwipeableArticleCard from './SwipeableArticleCard'

interface Props {
  articles: Article[]
  basePath: string
  emptyMessage: string
  hasMore: boolean
  fetchMoreArticles: () => Promise<void>
  refetch: () => Promise<any>
  filter?: (a: Article) => boolean
}

export default (props: Props) => {
  const {
    basePath,
    fetchMoreArticles,
    refetch,
    hasMore,
    filter = () => true,
    emptyMessage = 'No more article to read'
  } = props
  const ref = createRef<HTMLUListElement>()

  const isFetching = useInfiniteScroll(ref, fetchMoreArticles)
  const isMobileDisplay = useMedia('(max-width: 767px)')

  const articles = props.articles.filter(filter)

  if (articles.length === 0) {
    if (hasMore) {
      refetch()
    } else {
      return (
        <Empty>
          <ButtonIcon title="Refresh" icon="refresh" onClick={() => refetch()} />
          <br />
          <span>{emptyMessage}</span>
        </Empty>
      )
    }
  }

  return (
    <ul className={styles.list} ref={ref}>
      {articles.map(article => (
        <li key={`article-${article.id}`}>
          {isMobileDisplay && !article.isOffline ? (
            <SwipeableArticleCard article={article} readMoreBasePath={basePath} />
          ) : (
            <ArticleCard article={article} readMoreBasePath={basePath} />
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
