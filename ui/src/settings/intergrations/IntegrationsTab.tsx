import React from 'react'

import Panel from '../../components/Panel'
import { usePageTitle } from '../../hooks'
import IncomingWebhookSection from './incoming-webhook/IncomingWebhookSection'

export default () => {
  usePageTitle('Settings - Integrations')

  return (
    <Panel>
      <IncomingWebhookSection />
    </Panel>
  )
}
