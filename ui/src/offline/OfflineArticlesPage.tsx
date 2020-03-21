import React, { useCallback, useEffect, useContext } from 'react'
import { RouteComponentProps } from 'react-router'

import ArticleList from '../articles/components/ArticleList'
import { DisplayMode } from '../articles/components/ArticlesDisplayMode'
import ArticlesPageMenu from '../articles/components/ArticlesPageMenu'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import ErrorPanel from '../error/ErrorPanel'
import { getURLParam, matchState } from '../helpers'
import Page from '../layout/Page'
import { GetArticlesQuery, GetArticlesResult } from './dao/articles'
import { LocalConfigurationContext } from '../context/LocalConfigurationContext'

type AllProps = OfflineProps & RouteComponentProps

export const OfflineArticlesPage = ({ offlineArticles, fetchOfflineArticles, location }: AllProps) => {
  const { localConfiguration } = useContext(LocalConfigurationContext)

  const params = new URLSearchParams(location.search)
  const query: GetArticlesQuery = {
    limit: getURLParam<number>(params, 'limit', localConfiguration.limit),
    sortOrder: getURLParam<string>(params, 'sort', localConfiguration.sortOrders.offline)
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
    <Page title={title} actions={<ArticlesPageMenu refresh={refetch} mode={DisplayMode.offline} req={query} />}>
      {render(data, error, loading)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
