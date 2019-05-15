import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router-dom'

import { Category } from '../categories/models'
import ButtonIcon from '../common/ButtonIcon'
import Center from '../common/Center'
import { matchResponse } from '../common/helpers'
import Loader from '../common/Loader'
import Page from '../common/Page'
import Panel from '../common/Panel'
import ErrorPanel from '../error/ErrorPanel'
import useKeyboard from '../hooks/useKeyboard'
import ArticleContent from './components/ArticleContent'
import ArticleHeader from './components/ArticleHeader'
import ArticleMenu from './components/ArticleMenu'
import MarkAsButton from './components/MarkAsButton'
import { GetArticleResponse } from './models'
import { GetArticle } from './queries'

interface Props {
  category?: Category
}

type AllProps = Props & RouteComponentProps<{ id: string }>

export default ({ category, match, history }: AllProps) => {
  const { id } = match.params

  let title = 'Articles to read'
  if (category) {
    title = category.title
  }
  if (match.path === '/history/:id') {
    title = 'History'
  }

  // useKeyboard('backspace', () => history.push(redirect))
  useKeyboard('backspace', () => history.goBack())

  const { data, error, loading } = useQuery<GetArticleResponse>(GetArticle, {
    variables: { id }
  })

  const render = matchResponse<GetArticleResponse>({
    Loading: () => (
      <Center>
        <Loader />
      </Center>
    ),
    Error: err => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ article }) => {
      if (article) {
        article.isOffline = false
        return (
          <>
            <ArticleHeader article={article}>
              <ArticleMenu article={article} keyboard />
            </ArticleHeader>
            <ArticleContent article={article} />
            <MarkAsButton article={article} floating onSuccess={() => history.goBack()} keyboard />
          </>
        )
      }
      return <ErrorPanel title="Not found">Article #{id} not found.</ErrorPanel>
    },
    Other: () => <p>OTHER</p>
  })

  return (
    <Page
      title={title}
      subtitle={data && data.article ? data.article.title : ''}
      actions={<ButtonIcon onClick={history.goBack} icon="arrow_back" title="back to the list" />}
    >
      <Panel style={{ flex: '1 1 auto' }}>{render(data, error, loading)}</Panel>
    </Page>
  )
}
