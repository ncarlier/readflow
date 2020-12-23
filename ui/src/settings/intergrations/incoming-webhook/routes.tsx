import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router'
import AddIncomingWebhookForm from './AddIncomingWebhookForm'
import EditIncomingWebhookTab from './EditIncomingWebhookTab'

export default () => {
  const { path } = useRouteMatch()
  return (
    <Switch>
      <Route exact path={path + '/add'} component={AddIncomingWebhookForm} />
      <Route exact path={path + '/:id'} component={EditIncomingWebhookTab} />
    </Switch>
  )
}
