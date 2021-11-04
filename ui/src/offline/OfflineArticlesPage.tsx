import React, { useCallback, useEffect, useState } from 'react'
import { useLocation } from 'react-router'

import { ArticleList, ArticlesPageMenu } from '../articles/components'
import { GetArticlesRequest, GetArticlesResponse } from '../articles/models'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { LocalConfiguration, SortOrder, useLocalConfiguration } from '../contexts/LocalConfigurationContext'
import { Center, ErrorPanel, Loader, Panel } from '../components'
import { getURLParam, matchState } from '../helpers'
import { Page } from '../layout'

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

const OfflineArticlesPage = ({ offlineArticles, fetchOfflineArticles }: OfflineProps) => {
  const { localConfiguration } = useLocalConfiguration()
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
          variant={localConfiguration.display.offline.mode}
        />
      )
    },
  })

  const nb = data ? data.articles.totalCount : 0
  const title = `${nb} offline article${nb > 0 ? 's' : ''}`

  return (
    <Page title={title} actions={<ArticlesPageMenu refresh={refetch} variant="offline" req={req} />}>
      {render(loading, data, error)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
