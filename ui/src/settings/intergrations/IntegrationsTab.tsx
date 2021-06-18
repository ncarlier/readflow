import React from 'react'

import Panel from '../../components/Panel'
import { usePageTitle } from '../../hooks'
import FeedpushrSection from './feedpushr/FeedpushrSection'
import IncomingWebhookSection from './incoming-webhook/IncomingWebhookSection'
import OutgoingWebhookSection from './outgoing-webhook/OutgoingWebhookSection'

export default () => {
  usePageTitle('Settings - Integrations')

  return (
    <Panel>
      <FeedpushrSection />
      <IncomingWebhookSection />
      <OutgoingWebhookSection />
    </Panel>
  )
}
