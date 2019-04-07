import * as React from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import ArticlesRoutes from './articles/Routes'
import OfflineRoutes from './offline/Routes'
import CategoryRoutes from './categories/Routes'
import AboutPage from './about/AboutPage'
import SettingsPage from './settings/SettingsPage'
import ErrorPage from './error/ErrorPage'

const Routes = () => (
  <Switch>
    <Redirect exact from="/" to="/unread" />
    <Route path="/unread" component={ArticlesRoutes} />
    <Route path="/history" component={ArticlesRoutes} />
    <Route path="/offline" component={OfflineRoutes} />
    <Route path="/categories/:id" component={CategoryRoutes} />
    <Route path="/settings" component={SettingsPage} />
    <Route path="/about" component={AboutPage} />
    <Route component={() => <ErrorPage title="Not Found">Page not found</ErrorPage>} />
  </Switch>
)

export default Routes
