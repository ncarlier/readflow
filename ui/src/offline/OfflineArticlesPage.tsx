import React, { useEffect } from 'react'

import Page from  '../common/Page'
import ErrorPanel from '../error/ErrorPanel'
import { matchResponse } from '../common/helpers'
import Loader from '../common/Loader'
import Panel from '../common/Panel'
import ArticleList from '../articles/components/ArticleList'
import { Article } from '../articles/models'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { RouteComponentProps } from 'react-router'

type AllProps = OfflineProps & RouteComponentProps

export const OfflineArticlesPage = ({offlineArticles, fetchOfflineArticles, match}: AllProps) => {
  const { data, error, loading } = offlineArticles

  useEffect(() => {
    fetchOfflineArticles()
  }, [])
  
  const render = matchResponse<Article[]>({
    Loading: () => <Loader />,
    Error: (err) => <Panel><ErrorPanel>{err.message}</ErrorPanel></Panel>,
    Data: (d) => <ArticleList
      articles={d}
      basePath={match.path}
      emptyMessage="No offline articles"
      fetchMoreArticles={ () => console.log('TODO: fetchMoreArticles') }
    />,
    Other: () => <Panel><ErrorPanel>Unable to fetch articles!</ErrorPanel></Panel>
  })

  return (
    <Page title="Offline articles">
      {render(data, error, loading)}
    </Page>
  )
}

export default connectOffline(OfflineArticlesPage)
