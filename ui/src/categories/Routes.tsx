import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'

import ArticlePage from '../articles/ArticlePage'
import ArticlesPage from '../articles/ArticlesPage'
import { matchResponse } from '../common/helpers'
import Loader from '../common/Loader'
import Page from '../common/Page'
import ErrorPage from '../error/ErrorPage'
import { ConnectedReduxProps } from '../store'
import { GetCategoryResponse } from './models'
import { GetCategory } from './queries'

type AllProps = RouteComponentProps<{ id: string }> & ConnectedReduxProps

export default ({ match }: AllProps) => {
  const { id } = match.params
  const { data, error, loading } = useQuery<GetCategoryResponse>(GetCategory, {
    variables: { id }
  })

  const render = matchResponse<GetCategoryResponse>({
    Loading: () => (
      <Page>
        <Loader />
      </Page>
    ),
    Error: err => <ErrorPage>{err.message}</ErrorPage>,
    Data: ({ category }) => {
      if (category) {
        return (
          <Switch>
            <Route
              exact
              path={match.path + '/'}
              component={(props: RouteComponentProps) => <ArticlesPage category={category} {...props} />}
            />
            <Route
              path={match.path + '/articles/:id'}
              component={(props: RouteComponentProps<{ id: string }>) => <ArticlePage category={category} {...props} />}
            />
          </Switch>
        )
      } else {
        return <ErrorPage title="Not found">Category #${id} not found.</ErrorPage>
      }
    },
    Other: () => <ErrorPage>Unable to fetch category #${id} details!</ErrorPage>
  })

  return <>{render(data, error, loading)}</>
}
