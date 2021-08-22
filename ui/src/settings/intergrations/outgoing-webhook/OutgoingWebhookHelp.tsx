import React from 'react'

import { HelpLink, Logo } from '../../../components'
import HelpSection from '../../HelpSection'

export default () => (
  <HelpSection>
    <Logo name="webhook" style={{ maxWidth: '10%', verticalAlign: 'middle' }} />
    <span>
      Use outgoing webhooks to send articles to external services. <br />
    </span>
    <HelpLink href="https://docs.readflow.app/en/integrations/outgoing-webhook/">Help</HelpLink>
  </HelpSection>
)
