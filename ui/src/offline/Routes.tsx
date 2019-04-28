import React from 'react'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'

import OfflineArticlePage from './OfflineArticlePage'
import OfflineArticlesPage from './OfflineArticlesPage'

type AllProps = RouteComponentProps<{}>

export default ({ match }: AllProps) => (
  <Switch>
    <Route exact path={match.path + '/'} component={OfflineArticlesPage} />
    <Route exact path={match.path + '/:id'} component={OfflineArticlePage} />
  </Switch>
)
