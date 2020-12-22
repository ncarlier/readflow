import React from 'react'
import { useQuery } from '@apollo/client'
import { Route, RouteComponentProps, Switch, useRouteMatch } from 'react-router-dom'

import ArticlePage from '../articles/ArticlePage'
import ArticlesPage from '../articles/ArticlesPage'
import Loader from '../components/Loader'
import ErrorPage from '../error/ErrorPage'
import { matchResponse } from '../helpers'
import Page from '../layout/Page'
import { GetCategoryResponse } from './models'
import { GetCategory } from './queries'
import AddArticlePage from '../articles/AddArticlePage'

export default () => {
  const match = useRouteMatch<{ id: string }>()

  const { id } = match.params
  const { data, error, loading } = useQuery<GetCategoryResponse>(GetCategory, {
    variables: { id },
  })

  const render = matchResponse<GetCategoryResponse>({
    Loading: () => (
      <Page>
        <Loader />
      </Page>
    ),
    Error: (err) => <ErrorPage>{err.message}</ErrorPage>,
    Data: ({ category }) => {
      if (category) {
        return (
          <Switch>
            <Route
              exact
              path={match.path + '/'}
              component={(props: RouteComponentProps) => (
                <ArticlesPage variant="unread" category={category} {...props} />
              )}
            />
            <Route
              exact
              path={match.path + '/add'}
              component={(props: RouteComponentProps) => <AddArticlePage category={category} {...props} />}
            />
            <Route
              path={match.path + '/:id'}
              component={(props: RouteComponentProps<{ id: string }>) => (
                <ArticlePage title={category.title} {...props} />
              )}
            />
          </Switch>
        )
      } else {
        return <ErrorPage title="Not found">Category #${id} not found.</ErrorPage>
      }
    },
  })

  return <>{render(loading, data, error)}</>
}
