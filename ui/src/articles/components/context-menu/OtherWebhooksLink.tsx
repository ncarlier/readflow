import React, { useState } from 'react'

import LinkIcon from '../../../components/LinkIcon'
import { OutgoingWebhook } from '../../../settings/intergrations/outgoing-webhook/models'
import OutgoingWebhooksPanel from './OutgoingWebhooksPanel'
import Overlay from '../../../components/Overlay'

interface Props {
  webhooks: OutgoingWebhook[]
  sendArticle: (alias: string) => any
}

const OtherWebhooksLink = ({ webhooks, sendArticle }: Props) => {
  const [isVisible, setIsVisible] = useState(false)
  const showOverlay = () => setIsVisible(true)
  const hideOverlay = () => setIsVisible(false)

  if (webhooks.length > 1) {
    return (
      <>
        <LinkIcon title="Send article to ..." icon="backup" onClick={showOverlay}>
          <span>Send to ...</span>
        </LinkIcon>
        <Overlay visible={isVisible}>
          <OutgoingWebhooksPanel onCancel={hideOverlay} sendArticle={sendArticle} webhooks={webhooks} />
        </Overlay>
      </>
    )
  }
  return null
}

export default OtherWebhooksLink
