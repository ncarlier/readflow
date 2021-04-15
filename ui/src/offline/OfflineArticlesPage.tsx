import React, { useCallback, useContext, useEffect, useState } from 'react'
import { useLocation } from 'react-router'

import ArticleList from '../articles/components/ArticleList'
import ArticlesPageMenu from '../articles/components/ArticlesPageMenu'
import { GetArticlesRequest, GetArticlesResponse } from '../articles/models'
import Center from '../components/Center'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { LocalConfiguration, LocalConfigurationContext, SortOrder } from '../context/LocalConfigurationContext'
import ErrorPanel from '../error/ErrorPanel'
import { getURLParam, matchState } from '../helpers'
import Page from '../layout/Page'

const emptyQuery = {
  limit: null,
  sortOrder: null,
  afterCursor: null,
  category: null,
  starred: null,
  status: null,
}

const buildQueryFromLocation = (qs: string, localConfiguration: LocalConfiguration): GetArticlesRequest => {
  const params = new URLSearchParams(qs)
  return {
    ...emptyQuery,
    limit: getURLParam<number>(params, 'limit', localConfiguration.limit),
    sortBy: null,
    sortOrder: getURLParam<SortOrder>(params, 'sort', localConfiguration.display.offline.order),
    query: null,
  }
}

export const OfflineArticlesPage = ({ offlineArticles, fetchOfflineArticles }: OfflineProps) => {
  const { localConfiguration } = useContext(LocalConfigurationContext)
  const location = useLocation()
  const [req, setReq] = useState<GetArticlesRequest>(buildQueryFromLocation(location.search, localConfiguration))
  const { data, error, loading } = offlineArticles

  const refetch = useCallback(async () => {
    fetchOfflineArticles(req)
  }, [fetchOfflineArticles, req])

  const fetchMoreArticles = useCallback(async () => {
    if (!loading && data && data.articles.hasNext) {
      fetchOfflineArticles({ ...req, afterCursor: data.articles.endCursor })
    }
  }, [loading, data, req, fetchOfflineArticles])

  useEffect(() => {
    setReq(buildQueryFromLocation(location.search, localConfiguration))
  }, [location.search, localConfiguration])

  useEffect(() => {
    fetchOfflineArticles(req)
  }, [fetchOfflineArticles, req])

  const render = matchState<GetArticlesResponse>({
    Loading: () => (
      <Center>
        <Loader />
      </Center>
    ),
    Error: (err) => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: (d) => {
      return (
        <ArticleList
          articles={d.articles.entries}
          emptyMessage="no offline articles"
          hasMore={d.articles.hasNext}
          refetch={refetch}
          fetchMoreArticles={fetchMoreArticles}
        />
      )
    },
  })

  let title = ' '
  if (data) {
    const { totalCount } = data.articles
    const plural = totalCount > 1 ? ' articles' : ' article'
    title = totalCount + ' offline ' + plural
  }
  return (
    <Page title={title} actions={<ArticlesPageMenu refresh={refetch} variant="offline" req={req} />}>
      {render(loading, data, error)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
