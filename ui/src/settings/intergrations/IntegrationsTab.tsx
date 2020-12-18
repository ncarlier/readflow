import React from 'react'

import Panel from '../../components/Panel'
import { usePageTitle } from '../../hooks'
import IncomingWebhookSection from './incoming-webhook/IncomingWebhookSection'
import OutgoingWebhookSection from './outgoing-webhook/OutgoingWebhookSection'

export default () => {
  usePageTitle('Settings - Integrations')

  return (
    <Panel>
      <IncomingWebhookSection />
      <OutgoingWebhookSection />
    </Panel>
  )
}
