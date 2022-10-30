import React, { useCallback, useState } from 'react'

import { Kbd, LinkIcon, Loader, Overlay } from '../../../components'
import { OutgoingWebhook } from '../../../settings/intergrations/outgoing-webhook/models'
import OutgoingWebhooksPanel from './OutgoingWebhooksPanel'

interface Props {
  webhooks: OutgoingWebhook[]
  sendArticle: (alias: string) => Promise<any>
  keyboard?: boolean
}

const OutgoingWebhooksLink = ({ webhooks, sendArticle, keyboard }: Props) => {
  const [loading, setLoading] = useState(false)
  const [isVisible, setIsVisible] = useState(false)
  const showOverlay = () => setIsVisible(true)
  const hideOverlay = () => setIsVisible(false)
  
  const handleSendArticle = useCallback(async () => {
    const defaultWebhook = webhooks.find((webhook) => webhook.is_default)
    if (!defaultWebhook) {
      return
    }
    setLoading(true)
    try {
      await sendArticle(defaultWebhook.alias)
    } finally {
      setLoading(false)
    }
  }, [webhooks, sendArticle])

  if (loading) {
    return <Loader center />
  }
  if (webhooks.length > 0) {
    return (
      <>
        <LinkIcon title="Send article to ..." icon="backup" onClick={showOverlay}>
          <span>Send to ...</span>
          {keyboard && <Kbd keys="shift+s" onKeypress={handleSendArticle} />}
        </LinkIcon>
        <Overlay visible={isVisible}>
          <OutgoingWebhooksPanel onCancel={hideOverlay} sendArticle={sendArticle} webhooks={webhooks} />
        </Overlay>
      </>
    )
  }
  return null
}

export default OutgoingWebhooksLink
