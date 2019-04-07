import React from 'react'
import { RouteComponentProps, Route, Switch } from 'react-router-dom'

import ArticlesPage from './ArticlesPage'
import ArticlePage from './ArticlePage'

import { ConnectedReduxProps } from '../store'

// Combine both state + dispatch props - as well as any props we want to pass - in a union type.
type AllProps = RouteComponentProps<{}> & ConnectedReduxProps

export default ({match}: AllProps) => (
  <Switch>
    <Route exact path={match.path + '/'} component={ArticlesPage} />
    <Route path={match.path + '/:id'} component={ArticlePage} />
  </Switch>
)
