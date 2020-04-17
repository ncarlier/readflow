import React, { useCallback, useContext, useEffect } from 'react'
import { RouteComponentProps } from 'react-router'

import ArticleList from '../articles/components/ArticleList'
import { DisplayMode } from '../articles/components/ArticlesDisplayMode'
import ArticlesPageMenu from '../articles/components/ArticlesPageMenu'
import { GetArticlesRequest, GetArticlesResponse } from '../articles/models'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { LocalConfigurationContext } from '../context/LocalConfigurationContext'
import ErrorPanel from '../error/ErrorPanel'
import { getURLParam, matchState } from '../helpers'
import Page from '../layout/Page'

type AllProps = OfflineProps & RouteComponentProps

export const OfflineArticlesPage = ({ offlineArticles, fetchOfflineArticles, location }: AllProps) => {
  const { localConfiguration } = useContext(LocalConfigurationContext)

  const params = new URLSearchParams(location.search)
  const query: GetArticlesRequest = {
    limit: getURLParam<number>(params, 'limit', localConfiguration.limit),
    sortOrder: getURLParam<string>(params, 'sort', localConfiguration.sortOrders.offline),
    afterCursor: null,
    category: null,
    starred: null,
    status: null
  }

  const { data, error, loading } = offlineArticles

  useEffect(() => {
    console.log(query)
    fetchOfflineArticles(query)
  }, [location])

  const refetch = useCallback(async () => {
    fetchOfflineArticles(query)
  }, [data])

  const fetchMoreArticles = useCallback(async () => {
    if (!loading && data && data.articles.hasNext) {
      fetchOfflineArticles({ ...query, afterCursor: data.articles.endCursor })
    }
  }, [data])

  const render = matchState<GetArticlesResponse>({
    Loading: () => <Loader />,
    Error: err => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: d => (
      <ArticleList
        articles={d.articles.entries}
        emptyMessage="no offline articles"
        hasMore={d.articles.hasNext}
        refetch={refetch}
        fetchMoreArticles={fetchMoreArticles}
      />
    )
  })

  let title = ' '
  if (data) {
    const { totalCount } = data.articles
    const plural = totalCount > 1 ? ' articles' : ' article'
    title = totalCount + ' offline ' + plural
  }

  return (
    <Page title={title} actions={<ArticlesPageMenu refresh={refetch} mode={DisplayMode.offline} req={query} />}>
      {render(data, error, loading)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
