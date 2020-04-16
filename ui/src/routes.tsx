import * as React from 'react'
import { Redirect, Route, Switch } from 'react-router-dom'

import ArticlesRoutes from './articles/Routes'
import CategoryRoutes from './categories/Routes'
import ErrorPage from './error/ErrorPage'
import GraphiQLPage from './graphiql/GraphiQLPage'
import OfflineRoutes from './offline/Routes'
import SettingsPage from './settings/SettingsPage'

const Routes = () => (
  <Switch>
    <Redirect exact from="/" to="/unread" />
    <Redirect exact from="/login" to="/unread" />
    <Route path="/unread" component={ArticlesRoutes} />
    <Route path="/history" component={ArticlesRoutes} />
    <Route path="/starred" component={ArticlesRoutes} />
    <Route path="/offline" component={OfflineRoutes} />
    <Route path="/categories/:id" component={CategoryRoutes} />
    <Route path="/settings" component={SettingsPage} />
    <Route path="/graphiql" component={GraphiQLPage} />
    <Route component={() => <ErrorPage title="Not Found">Page not found</ErrorPage>} />
  </Switch>
)

export default Routes
