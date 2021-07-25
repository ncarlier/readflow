import React, { useCallback, useState } from 'react'

import Loader from '../../../components/Loader'
import { OutgoingWebhook } from '../../../settings/intergrations/outgoing-webhook/models'
import Button from '../../../components/Button'
import Panel from '../../../components/Panel'
import Logo from '../../../logos/Logo'

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

  return (
    <Panel>
      {loading && <Loader blur />}
      <header>
        <h1>Send article to...</h1>
      </header>
      <section>
        {webhooks.map((webhook) => (
          <Button
            key={`wh${webhook.id}`}
            variant="flat"
            style={{ width: '10rem', padding: '1rem' }}
            title={`Send article to ${webhook.alias}`}
            onClick={() => handleSendArticle(webhook.alias).then(onCancel)}
          >
            <Logo name={webhook.provider} style={{ maxWidth: '2em', verticalAlign: 'middle' }} />
            <p>Send to {webhook.alias}</p>
          </Button>
        ))}
      </section>
      <footer>
        <Button title="Cancel" onClick={onCancel}>
          Cancel
        </Button>
      </footer>
    </Panel>
  )
}
