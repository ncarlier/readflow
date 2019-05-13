import React, { useCallback, useEffect } from 'react'
import { RouteComponentProps } from 'react-router'

import ArticleList from '../articles/components/ArticleList'
import { DisplayMode } from '../articles/components/ArticlesDisplayMode'
import ArticlesPageMenu from '../articles/components/ArticlesPageMenu'
import { getURLParam, matchState } from '../common/helpers'
import Loader from '../common/Loader'
import Page from '../common/Page'
import Panel from '../common/Panel'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import ErrorPanel from '../error/ErrorPanel'
import { GetArticlesQuery, GetArticlesResult } from './dao/articles'

type AllProps = OfflineProps & RouteComponentProps

export const OfflineArticlesPage = ({ offlineArticles, fetchOfflineArticles, match, location }: AllProps) => {
  const params = new URLSearchParams(location.search)
  const query: GetArticlesQuery = {
    limit: getURLParam<number>(params, 'limit', 10),
    sortOrder: getURLParam<string>(params, 'sort', 'asc')
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
      fetchOfflineArticles({ ...query, afterCursor: data.endCursor })
    }
  }, [data])

  const render = matchState<GetArticlesResult>({
    Loading: () => <Loader />,
    Error: err => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: d => (
      <ArticleList
        articles={d.entries}
        basePath={match.path}
        emptyMessage="no offline articles"
        hasMore={d.hasNext}
        refetch={refetch}
        fetchMoreArticles={fetchMoreArticles}
      />
    )
  })

  let title = ' '
  if (data) {
    const { totalCount } = data
    const plural = totalCount > 1 ? ' articles' : ' article'
    title = data.totalCount + ' offline ' + plural
  }

  return (
    <Page title={title} actions={<ArticlesPageMenu refresh={refetch} mode={DisplayMode.offline} />}>
      {render(data, error, loading)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
