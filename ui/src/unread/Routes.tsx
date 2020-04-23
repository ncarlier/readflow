import React from 'react'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'

import ArticlePage from '../articles/ArticlePage'
import ArticlesPage from '../articles/ArticlesPage'
import AddArticlePage from '../articles/AddArticlePage'

export default ({ match }: RouteComponentProps) => (
  <Switch>
    <Route
      exact
      path={match.path + '/'}
      component={(props: RouteComponentProps) => <ArticlesPage variant="unread" {...props} />}
    />
    <Route exact path={match.path + '/add'} component={AddArticlePage} />
    <Route
      path={match.path + '/:id'}
      component={(props: RouteComponentProps<{ id: string }>) => <ArticlePage title="Articles to read" {...props} />}
    />
  </Switch>
)
