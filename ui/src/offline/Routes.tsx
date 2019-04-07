import React from 'react'
import { RouteComponentProps, Route, Switch } from 'react-router-dom'

import OfflineArticlesPage from './OfflineArticlesPage'
import OfflineArticlePage from './OfflineArticlePage'

type AllProps = RouteComponentProps<{}>

export default ({match}: AllProps) => (
  <Switch>
    <Route exact path={match.path + '/'} component={OfflineArticlesPage} />
    <Route exact path={match.path + '/:id'} component={OfflineArticlePage} />
  </Switch>
)
