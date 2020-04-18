import React from 'react'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'

import ArticlePage from '../articles/ArticlePage'
import ArticlesPage from '../articles/ArticlesPage'

export default ({ match }: RouteComponentProps) => (
  <Switch>
    <Route
      exact
      path={match.path + '/'}
      component={(props: RouteComponentProps) => <ArticlesPage variant="starred" {...props} />}
    />
    <Route
      path={match.path + '/:id'}
      component={(props: RouteComponentProps<{ id: string }>) => <ArticlePage title="Starred articles" {...props} />}
    />
  </Switch>
)
