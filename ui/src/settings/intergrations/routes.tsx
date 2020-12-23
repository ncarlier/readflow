import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router'
import IntegrationsTab from './IntegrationsTab'
import IncomingWebhookRoutes from './incoming-webhook/routes'
import OutgoingWebhookRoutes from './outgoing-webhook/routes'

export default () => {
  const { path } = useRouteMatch()
  return (
    <Switch>
      <Route exact path={path + '/'} component={IntegrationsTab} />
      <Route path={path + '/incoming-webhooks'} component={IncomingWebhookRoutes} />
      <Route path={path + '/outgoing-webhooks'} component={OutgoingWebhookRoutes} />
    </Switch>
  )
}
