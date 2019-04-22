import React from 'react'
import { RouteComponentProps, Route, Switch } from 'react-router-dom'

import { ConnectedReduxProps } from '../store'

import ArticlesPage from './ArticlesPage'
import ArticlePage from './ArticlePage'
import AddArticlePage from './AddArticlePage'

// Combine both state + dispatch props - as well as any props we want to pass - in a union type.
type AllProps = RouteComponentProps<{}> & ConnectedReduxProps

export default ({match}: AllProps) => (
  <Switch>
    <Route exact path={match.path + '/'} component={ArticlesPage} />
    <Route exact path={match.path + '/add'} component={AddArticlePage} />
    <Route path={match.path + '/:id'} component={ArticlePage} />
  </Switch>
)
