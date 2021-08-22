import React, { useEffect } from 'react'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'

import { ArticleContent, ArticleHeader, ArticleContextMenu } from '../articles/components'
import { Article } from '../articles/models'
import { ButtonIcon, Center, ErrorPanel, Loader, Panel } from '../components'
import { connectOffline, OfflineProps } from '../containers/OfflineContainer'
import { matchResponse } from '../helpers'
import { Page } from '../layout'

type AllProps = RouteComponentProps<{ id: string }> & OfflineProps

const Actions = () => <ButtonIcon as={Link} to="/offline" icon="arrow_back" title="back to the list" />

const OfflineArticlePage = ({ match, offlineArticles, fetchOfflineArticle }: AllProps) => {
  const { id } = match.params

  const { selected: data, error, loading } = offlineArticles

  useEffect(() => {
    fetchOfflineArticle(parseInt(id, 10))
  }, [fetchOfflineArticle, id])

  const render = matchResponse<Article>({
    Loading: () => (
      <Center>
        <Loader />
      </Center>
    ),
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: (article) => {
      if (article) {
        article.isOffline = true
        return (
          <>
            <ArticleHeader article={article}>
              <ArticleContextMenu article={article} />
            </ArticleHeader>
            <ArticleContent article={article} />
          </>
        )
      }
      return <ErrorPanel title="Not found">Article #{id} not found.</ErrorPanel>
    },
  })

  return (
    <Page title="Offline articles" subtitle={data && data.title} actions={<Actions />}>
      <Panel style={{ flex: '1 1 auto' }}>{render(loading, data, error)}</Panel>
    </Page>
  )
}

export default connectOffline(OfflineArticlePage)
