import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'

import Kbd from '../../components/Kbd'
import LinkIcon from '../../components/LinkIcon'
import Loader from '../../components/Loader'
import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'
import { OutgoingWebhook } from '../../settings/intergrations/outgoing-webhook/models'
import { Article } from '../models'
import { SendArticleToOutgoingWebhook } from '../queries'

interface SendArticleFields {
  id: number
  alias: string
  noShortcuts?: boolean
}

interface Props {
  webhook: OutgoingWebhook
  article: Article
  keyboard?: boolean
}

export default (props: Props) => {
  const [loading, setLoading] = useState(false)
  const { article, webhook, keyboard = false } = props
  const { showMessage, showErrorMessage } = useContext(MessageContext)

  const [SendArticleToOutgoingWebhookMutation] = useMutation<SendArticleFields>(SendArticleToOutgoingWebhook)

  const sendArticle = useCallback(async () => {
    setLoading(true)
    try {
      const { alias } = webhook
      await SendArticleToOutgoingWebhookMutation({
        variables: { id: article.id, alias },
      })
      showMessage(`Article sent to ${alias}: ${article.title}`)
    } catch (err) {
      showErrorMessage(getGQLError(err))
    } finally {
      setLoading(false)
    }
  }, [SendArticleToOutgoingWebhookMutation, article, webhook, showMessage, showErrorMessage])

  if (loading) {
    return <Loader />
  }

  return (
    <LinkIcon title={`Send article to ${webhook.alias}`} icon="backup" onClick={sendArticle}>
      <span>Send to {webhook.alias}</span>
      {keyboard && webhook.is_default && <Kbd keys="shift+s" onKeypress={sendArticle} />}
    </LinkIcon>
  )
}
