import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router-dom'

import { Category } from '../categories/models'
import ButtonIcon from '../components/ButtonIcon'
import Center from '../components/Center'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import ErrorPanel from '../error/ErrorPanel'
import { matchResponse } from '../helpers'
import useKeyboard from '../hooks/useKeyboard'
import Page from '../layout/Page'
import ArticleContent from './components/ArticleContent'
import ArticleHeader from './components/ArticleHeader'
import ArticleMenu from './components/ArticleMenu'
import MarkAsButton from './components/MarkAsButton'
import StarsButton from './components/StarsButton'
import { GetArticleResponse } from './models'
import { GetArticle } from './queries'

interface Props {
  title: string
  category?: Category
}

type AllProps = Props & RouteComponentProps<{ id: string }>

export default ({ title, match, history }: AllProps) => {
  const { id } = match.params

  const goBack = () => history.goBack()

  useKeyboard('backspace', goBack)

  const { data, error, loading } = useQuery<GetArticleResponse>(GetArticle, {
    variables: { id },
  })

  const render = matchResponse<GetArticleResponse>({
    Loading: () => (
      <Center>
        <Loader />
      </Center>
    ),
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ article }) => {
      if (article) {
        article.isOffline = false
        return (
          <>
            <ArticleHeader article={article}>
              <StarsButton article={article} keyboard />
              <ArticleMenu article={article} keyboard />
            </ArticleHeader>
            <ArticleContent article={article} />
            <MarkAsButton article={article} floating onSuccess={goBack} keyboard />
          </>
        )
      }
      return <ErrorPanel title="Not found">Article #{id} not found.</ErrorPanel>
    },
    Other: () => <p>OTHER</p>,
  })

  return (
    <Page
      title={title}
      subtitle={data && data.article ? data.article.title : ''}
      actions={<ButtonIcon onClick={goBack} icon="arrow_back" title="back to the list" />}
    >
      <Panel style={{ flex: '1 1 auto' }}>{render(data, error, loading)}</Panel>
    </Page>
  )
}
