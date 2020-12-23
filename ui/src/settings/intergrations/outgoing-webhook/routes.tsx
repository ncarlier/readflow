import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router'
import AddOutgoingWebhookForm from './AddOutgoingWebhookForm'
import EditOutgoingWebhookTab from './EditOutgoingWebhookTab'

export default () => {
  const { path } = useRouteMatch()
  return (
    <Switch>
      <Route exact path={path + '/add'} component={AddOutgoingWebhookForm} />
      <Route exact path={path + '/:id'} component={EditOutgoingWebhookTab} />
    </Switch>
  )
}
