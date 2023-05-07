import React, { useCallback } from 'react'
import { useMutation, useQuery } from '@apollo/client'

import { Loader } from '../../../components'
import { getGQLError, matchResponse } from '../../../helpers'
import { GetOutgoingWebhooksResponse } from '../../../settings/intergrations/outgoing-webhook/models'
import { GetOutgoingWebhooks } from '../../../settings/intergrations/outgoing-webhook/queries'
import { Article, SendArticleToOutgoingWebhookResponse } from '../../models'
import { SendArticleToOutgoingWebhook } from '../../queries'
import { useMessage } from '../../../contexts'
import OutgoingWebhooksLink from './OutgoingWebhooksLink'

interface SendArticleFields {
  id: number
  alias: string
  noShortcuts?: boolean
}

interface Props {
  article: Article
  keyboard?: boolean
}

export default ({ article, keyboard }: Props) => {
  const { showMessage, showErrorMessage } = useMessage()
  const [SendArticleToOutgoingWebhookMutation] = useMutation<SendArticleToOutgoingWebhookResponse, SendArticleFields>(SendArticleToOutgoingWebhook)
  const { data, error, loading } = useQuery<GetOutgoingWebhooksResponse>(GetOutgoingWebhooks)
  
  const sendArticle = useCallback(
    async (alias: string) => {
      try {
        const result = await SendArticleToOutgoingWebhookMutation({
          variables: { id: article.id, alias },
        })
        let message = `Article "${article.title}" sent to ${alias}`
        if (result.data) {
          const { url, text } = result.data.sendArticleToOutgoingWebhook
          if (text) {
            alert(text)
          }
          if (url) {
            message = `${message}\n[more](${url})`
          }
        }
        showMessage(message)
      } catch (err) {
        showErrorMessage(getGQLError(err))
      }
    },
    [SendArticleToOutgoingWebhookMutation, article, showMessage, showErrorMessage ]
  )

  const render = matchResponse<GetOutgoingWebhooksResponse>({
    Loading: () => (
      <li>
        <Loader center />
      </li>
    ),
    Error: (err) => <li>{err.message}</li>,
    Data: ({ outgoingWebhooks }) => (
      <li>
        <OutgoingWebhooksLink
          webhooks={outgoingWebhooks}
          keyboard={keyboard}
          sendArticle={sendArticle}
        />
      </li>
    ),
  })

  return <>{render(loading, data, error)}</>
}
