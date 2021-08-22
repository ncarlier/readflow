import React, { useCallback, useState } from 'react'

import { OutgoingWebhook } from '../../../settings/intergrations/outgoing-webhook/models'
import { Loader, Logo, LinkIcon } from '../../../components'

interface Props {
  webhooks: OutgoingWebhook[]
  sendArticle: (alias: string) => any
  onCancel: (e: any) => void
}

export default ({ webhooks, sendArticle, onCancel }: Props) => {
  const [loading, setLoading] = useState(false)
  const handleSendArticle = useCallback(
    async (alias: string) => {
      setLoading(true)
      try {
        await sendArticle(alias)
      } finally {
        setLoading(false)
      }
    },
    [sendArticle]
  )
  if (loading) {
    return <Loader blur />
  }

  return (
    <ul>
      {webhooks.map((webhook) => (
        <li key={`wh${webhook.id}`}>
          <LinkIcon
            icon={<Logo name={webhook.provider} style={{ maxWidth: '2em', verticalAlign: 'middle' }} />}
            title={`Send article to ${webhook.alias}`}
            onClick={() => handleSendArticle(webhook.alias).then(onCancel)}
          >
            <span>Send to {webhook.alias}</span>
          </LinkIcon>
        </li>
      ))}
      <li>
        <LinkIcon title="Cancel" icon="cancel" onClick={onCancel}>
          <span>Cancel</span>
        </LinkIcon>
      </li>
    </ul>
  )
}
