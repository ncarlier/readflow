import * as React from 'react'
import { Redirect, Route, Switch } from 'react-router-dom'

import CategoryRoutes from './categories/Routes'
import ErrorPage from './error/ErrorPage'
import GraphiQLPage from './graphiql/GraphiQLPage'
import HistoryRoutes from './history/Routes'
import OfflineRoutes from './offline/Routes'
import SettingsPage from './settings/SettingsPage'
import StarredRoutes from './starred/Routes'
import UnreadRoutes from './unread/Routes'

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
