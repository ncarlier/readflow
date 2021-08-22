import React from 'react'
import { useQuery } from '@apollo/client'
import { RouteComponentProps } from 'react-router-dom'

import { Category } from '../categories/models'
import { ButtonIcon, Center, ErrorPanel, Loader, Panel } from '../components'
import { matchResponse } from '../helpers'
import { useKeyboard } from '../hooks'
import { Page } from '../layout'
import { ArticleContent, ArticleHeader, MarkAsButton, StarsButton, ArticleContextMenu } from './components'
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
        return (
          <>
            <ArticleHeader article={article}>
              <StarsButton article={article} keyboard />
              <ArticleContextMenu article={article} keyboard />
            </ArticleHeader>
            <ArticleContent article={article} />
            <MarkAsButton article={article} floating onSuccess={goBack} keyboard />
          </>
        )
      }
      return <ErrorPanel title="Not found">Article #{id} not found.</ErrorPanel>
    },
  })

  return (
    <Page
      title={title}
      subtitle={data && data.article ? data.article.title : ''}
      actions={<ButtonIcon onClick={goBack} icon="arrow_back" title="back to the list" />}
    >
      <Panel style={{ flex: '1 1 auto', overflow: 'hidden' }}>{render(loading, data, error)}</Panel>
    </Page>
  )
}
