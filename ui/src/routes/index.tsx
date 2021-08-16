import * as React from 'react'
import { Redirect, Route, Switch } from 'react-router-dom'

import CategoryRoutes from './categories'
import ErrorPage from '../error/ErrorPage'
import GraphiQLPage from '../graphiql/GraphiQLPage'
import HistoryRoutes from './history'
import OfflineRoutes from './offline'
import SettingsPage from '../settings/SettingsPage'
import StarredRoutes from './starred'
import UnreadRoutes from './unread'

const Routes = () => (
  <Switch>
    <Redirect exact from="/" to="/unread" />
    <Redirect exact from="/login" to="/unread" />
    <Route path="/unread" component={UnreadRoutes} />
    <Route path="/history" component={HistoryRoutes} />
    <Route path="/starred" component={StarredRoutes} />
    <Route path="/offline" component={OfflineRoutes} />
    <Route path="/categories/:id" component={CategoryRoutes} />
    <Route path="/settings" component={SettingsPage} />
    <Route path="/graphiql" component={GraphiQLPage} />
    <Route component={() => <ErrorPage title="Not Found">Page not found</ErrorPage>} />
  </Switch>
)

export default Routes
