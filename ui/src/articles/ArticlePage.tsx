import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps, Redirect } from 'react-router-dom'

import { GetArticleResponse } from './models'
import { GetArticle } from './queries'
import Page from  '../common/Page'
import ArticleHeader from './components/ArticleHeader'
import ArticleContent from './components/ArticleContent'
import { matchResponse } from '../common/helpers'
import Loader from '../common/Loader'
import ErrorPanel from '../error/ErrorPanel'
import Panel from '../common/Panel'
import { Category } from '../categories/models'
import ButtonIcon from '../common/ButtonIcon'
import MarkAsButton from './components/MarkAsButton'
import ArticleMenu from './components/ArticleMenu'
import Center from '../common/Center';

type Props = {
  category?: Category
}

type AllProps = Props & RouteComponentProps<{id: string}>

export default ({ category, match }: AllProps) => {
  const { id } = match.params

  let title = 'Articles to read'
  let redirect = '/unread'
  if (category) {
    title = category.title
    redirect = `/categories/${category.id}`
  }
  if (match.path === '/history/:id') {
    title = 'History'
    redirect = '/history'
  }

  const { data, error, loading } = useQuery<GetArticleResponse>(GetArticle, {
    variables: {id}
  })
  
  const render = matchResponse<GetArticleResponse>({
    Loading: () => <Center><Loader /></Center>,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({article}) => 
      <>
        {article !== null ? 
          <>
            <ArticleHeader article={article}>
              <ArticleMenu article={article} />
            </ArticleHeader>
            <ArticleContent article={article} />
            <MarkAsButton article={article} floating />
          </>
          : <ErrorPanel title="Not found">Article #{id} not found.</ErrorPanel>
        }
      </>,
    Other: () => <Redirect to={redirect} />
  })

  return (
    <Page title={title}
          subtitle={data && data.article ? data.article.title : ''}
          actions={
            <ButtonIcon
              to={redirect} 
              icon="arrow_back"
              title="back to the list"
            />
          }>
      <Panel style={{flex: '1 1 auto'}}>
        {render(data, error, loading)}
      </Panel>
    </Page>
  )
}
