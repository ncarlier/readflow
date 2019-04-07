import React, { useEffect } from 'react'

import Page from  '../common/Page'
import ErrorPanel from '../error/ErrorPanel'
import { matchResponse } from '../common/helpers'
import Loader from '../common/Loader'
import Panel from '../common/Panel'
import { RouteComponentProps, Redirect } from 'react-router'
import ArticleHeader from '../articles/components/ArticleHeader'
import ArticleContent from '../articles/components/ArticleContent'
import { Article } from '../articles/models'
import { OfflineProps, connectOffline } from '../containers/OfflineContainer';

type AllProps = RouteComponentProps<{id: string}> & OfflineProps

export const OfflineArticlePage = ({match, offlineArticles, fetchOfflineArticle}: AllProps) => {
  const { id } = match.params
  
  const { selected: data, error, loading } = offlineArticles

  useEffect(() => {
    fetchOfflineArticle(parseInt(id, 10))
  }, [])
  
  const render = matchResponse<Article>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: (a) => <>
        {a !== null && (a.isOffline = true) ? 
          <>
            <ArticleHeader article={a} />
            <ArticleContent article={a} />
          </>
          : <ErrorPanel title="Not found">Article #${id} not found.</ErrorPanel>
        }
      </>,
    Other: () => <Redirect to="/offline" />
  })

  return (
    <Page title="Offline articles" subtitle={data && data.title}>
      <Panel style={{flex: '1 1 auto'}}>
        {render(data, error, loading)}
      </Panel>
    </Page>
  )
}

export default connectOffline(OfflineArticlePage)
