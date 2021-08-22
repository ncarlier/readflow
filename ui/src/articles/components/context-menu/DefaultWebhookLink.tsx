import React, { useCallback, useState } from 'react'

import { Kbd, LinkIcon, Loader } from '../../../components'
import { OutgoingWebhook } from '../../../settings/intergrations/outgoing-webhook/models'

interface Props {
  webhook?: OutgoingWebhook
  sendArticle: (alias: string) => Promise<any>
  keyboard?: boolean
}

const DefaultWebhookLink = ({ webhook, sendArticle, keyboard }: Props) => {
  const [loading, setLoading] = useState(false)
  const handleSendArticle = useCallback(async () => {
    if (!webhook) {
      return
    }
    setLoading(true)
    try {
      await sendArticle(webhook.alias)
    } finally {
      setLoading(false)
    }
  }, [webhook, sendArticle])

  if (!webhook) {
    return null
  }
  if (loading) {
    return <Loader center />
  }
  return (
    <LinkIcon title={`Send article to ${webhook.alias}`} icon="backup" onClick={handleSendArticle}>
      <span>Send to {webhook.alias}</span>
      {keyboard && webhook.is_default && <Kbd keys="shift+s" onKeypress={handleSendArticle} />}
    </LinkIcon>
  )
}

export default DefaultWebhookLink
