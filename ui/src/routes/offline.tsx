import React from 'react'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'

import OfflineArticlePage from '../offline/OfflineArticlePage'
import OfflineArticlesPage from '../offline/OfflineArticlesPage'

export default ({ match }: RouteComponentProps) => (
  <Switch>
    <Route exact path={match.path + '/'} component={OfflineArticlesPage} />
    <Route exact path={match.path + '/:id'} component={OfflineArticlePage} />
  </Switch>
)
