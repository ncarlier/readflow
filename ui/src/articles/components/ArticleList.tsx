import React, { useRef, createRef } from 'react'

import { Article } from '../models'
import ArticleCard from './ArticleCard'

import styles from './ArticleList.module.css'
import Empty from '../../common/Empty'
import useInfiniteScroll from '../../hooks/useInfiniteScroll'
import Panel from '../../common/Panel'
import Center from '../../common/Center'

type Props = {
  articles: Article[]
  basePath: string
  emptyMessage: string
  fetchMoreArticles: () => void
}

export default (props: Props) => {
  const {
    articles,
    basePath,
    fetchMoreArticles,
    emptyMessage = 'No more article to read'
  } = props
  const ref = createRef<HTMLUListElement>()
  const [isFetching, setIsFetching] = useInfiniteScroll(ref, onFetchMoreItems)
  
  if (articles.length === 0) {
    return <Empty>{ emptyMessage }</Empty>
  }

  async function onFetchMoreItems() {
    try {
      await fetchMoreArticles()
    } finally {
      setIsFetching(false)
    }
  }

  return (
    <ul className={styles.list} ref={ref}>
      {articles.map(article => (
        <li key={`article-${article.id}`}>
          <ArticleCard article={article} readMoreBasePath={basePath} />
        </li>
      ))}
      {isFetching && <li><Panel><Center>Fetching more articles...</Center></Panel></li>}
    </ul>
  )
}
