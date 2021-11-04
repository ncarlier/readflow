import * as React from 'react'
import { Redirect, Route, Switch } from 'react-router-dom'

import CategoryRoutes from './categories'
import { ErrorPage } from '../components/errors'
import GraphiQLPage from '../graphiql/GraphiQLPage'
import HistoryRoutes from './history'
import OfflineRoutes from './offline'
import SettingsPage from '../settings/SettingsPage'
import StarredRoutes from './starred'
import InboxRoutes from './inbox'
import ToReadRoutes from './to_read'

const Routes = () => (
  <Switch>
    <Redirect exact from="/" to="/inbox" />
    <Redirect exact from="/login" to="/inbox" />
    <Route path="/inbox" component={InboxRoutes} />
    <Route path="/to_read" component={ToReadRoutes} />
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
