import React, { useEffect, useCallback } from 'react'

import Page from  '../common/Page'
import ErrorPanel from '../error/ErrorPanel'
import { matchState, getURLParam } from '../common/helpers'
import Loader from '../common/Loader'
import Panel from '../common/Panel'
import ArticleList from '../articles/components/ArticleList'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { RouteComponentProps } from 'react-router'
import { GetArticlesQuery, GetArticlesResult } from './dao/articles'
import ArticlesPageMenu from '../articles/components/ArticlesPageMenu'

type AllProps = OfflineProps & RouteComponentProps

export const OfflineArticlesPage = ({offlineArticles, fetchOfflineArticles, match, location}: AllProps) => {
  const params = new URLSearchParams(location.search)
  const query: GetArticlesQuery = {
    limit: getURLParam<number>(params, 'limit', 10),
    sortOrder: getURLParam<string>(params, 'sort', 'asc'),
  }
  
  const { data, error, loading } = offlineArticles

  useEffect(() => {
    fetchOfflineArticles(query)
  }, [location])

  const refetch = useCallback(async () => {
    fetchOfflineArticles(query)
  }, [data])

  const fetchMoreArticles = useCallback(async () => {
    if (!loading && data && data.hasNext) {
      fetchOfflineArticles({...query, afterCursor: data.endCursor})
    }
  }, [data])
  
  const render = matchState<GetArticlesResult>({
    Loading: () => <Loader />,
    Error: (err) => <Panel><ErrorPanel>{err.message}</ErrorPanel></Panel>,
    Data: (d) => <ArticleList
      articles={d.entries}
      basePath={match.path}
      emptyMessage="no offline articles"
      hasMore={d.hasNext}
      refetch={refetch}
      fetchMoreArticles={ fetchMoreArticles }
    />,
  })

  let title = ' '
  if (data) {
    const {totalCount} = data
    const plural = totalCount > 1 ? " articles" : " article" 
    title = data.totalCount + ' offline ' + plural
  }

  return (
    <Page title={title} actions={
        <ArticlesPageMenu refresh={refetch} />
      }>
      {render(data, error, loading)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
