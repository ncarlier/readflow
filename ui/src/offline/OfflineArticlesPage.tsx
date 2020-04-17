import React, { useCallback, useContext, useEffect, useState } from 'react'
import { useLocation } from 'react-router'

import ArticleList from '../articles/components/ArticleList'
import { DisplayMode } from '../articles/components/ArticlesDisplayMode'
import ArticlesPageMenu from '../articles/components/ArticlesPageMenu'
import { GetArticlesRequest, GetArticlesResponse } from '../articles/models'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { LocalConfiguration, LocalConfigurationContext } from '../context/LocalConfigurationContext'
import ErrorPanel from '../error/ErrorPanel'
import { getURLParam, matchState } from '../helpers'
import Page from '../layout/Page'

const emptyQuery = {
  limit: null,
  sortOrder: null,
  afterCursor: null,
  category: null,
  starred: null,
  status: null
}

const buildQueryFromLocation = (qs: string, localConfiguration: LocalConfiguration): GetArticlesRequest => {
  const params = new URLSearchParams(qs)
  return {
    ...emptyQuery,
    limit: getURLParam<number>(params, 'limit', localConfiguration.limit),
    sortOrder: getURLParam<string>(params, 'sort', localConfiguration.sortOrders.offline)
  }
}

export const OfflineArticlesPage = ({ offlineArticles, fetchOfflineArticles }: OfflineProps) => {
  const { localConfiguration } = useContext(LocalConfigurationContext)
  const location = useLocation()
  const [query] = useState<GetArticlesRequest>(buildQueryFromLocation(location.search, localConfiguration))
  const { data, error, loading } = offlineArticles

  const refetch = useCallback(async () => {
    fetchOfflineArticles(query)
  }, [query, fetchOfflineArticles])

  useEffect(() => {
    if (data === undefined) {
      refetch()
    }
  }, [data, refetch])

  const fetchMoreArticles = useCallback(async () => {
    if (!loading && data && data.articles.hasNext) {
      fetchOfflineArticles({ ...query, afterCursor: data.articles.endCursor })
    }
  }, [loading, data, query, fetchOfflineArticles])

  const render = matchState<GetArticlesResponse>({
    Loading: () => <Loader />,
    Error: err => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: d => {
      return (
        <ArticleList
          articles={d.articles.entries}
          emptyMessage="no offline articles"
          hasMore={d.articles.hasNext}
          refetch={refetch}
          fetchMoreArticles={fetchMoreArticles}
        />
      )
    }
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
