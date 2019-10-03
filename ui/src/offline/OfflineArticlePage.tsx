import React, { useEffect } from 'react'
import { RouteComponentProps } from 'react-router'

import ArticleContent from '../articles/components/ArticleContent'
import ArticleHeader from '../articles/components/ArticleHeader'
import ArticleMenu from '../articles/components/ArticleMenu'
import { Article } from '../articles/models'
import ButtonIcon from '../components/ButtonIcon'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import ErrorPanel from '../error/ErrorPanel'
import { matchResponse } from '../helpers'
import Page from '../layout/Page'

type AllProps = RouteComponentProps<{ id: string }> & OfflineProps

const Actions = () => <ButtonIcon to="/offline" icon="arrow_back" title="back to the list" />

export const OfflineArticlePage = ({ match, offlineArticles, fetchOfflineArticle }: AllProps) => {
  const { id } = match.params

  const { selected: data, error, loading } = offlineArticles

  useEffect(() => {
    fetchOfflineArticle(parseInt(id, 10))
  }, [])

  const render = matchResponse<Article>({
    Loading: () => <Loader />,
    Error: err => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: article => {
      if (article) {
        article.isOffline = true
        return (
          <>
            <ArticleHeader article={article}>
              <ArticleMenu article={article} />
            </ArticleHeader>
            <ArticleContent article={article} />
          </>
        )
      }
      return <ErrorPanel title="Not found">Article #{id} not found.</ErrorPanel>
    },
    Other: () => <Loader />
  })

  return (
    <Page title="Offline articles" subtitle={data && data.title} actions={<Actions />}>
      <Panel style={{ flex: '1 1 auto' }}>{render(data, error, loading)}</Panel>
    </Page>
  )
}

export default connectOffline(OfflineArticlePage)
