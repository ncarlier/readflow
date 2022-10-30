import React, { useCallback } from 'react'
import { useMutation, useQuery } from '@apollo/client'

import { Loader } from '../../../components'
import { getGQLError, matchResponse } from '../../../helpers'
import { GetOutgoingWebhooksResponse } from '../../../settings/intergrations/outgoing-webhook/models'
import { GetOutgoingWebhooks } from '../../../settings/intergrations/outgoing-webhook/queries'
import { Article } from '../../models'
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
  const [SendArticleToOutgoingWebhookMutation] = useMutation<SendArticleFields>(SendArticleToOutgoingWebhook)
  const { data, error, loading } = useQuery<GetOutgoingWebhooksResponse>(GetOutgoingWebhooks)

  const sendArticle = useCallback(
    async (alias: string) => {
      try {
        await SendArticleToOutgoingWebhookMutation({
          variables: { id: article.id, alias },
        })
        showMessage(`Article sent to ${alias}: ${article.title}`)
      } catch (err) {
        showErrorMessage(getGQLError(err))
      }
    },
    [SendArticleToOutgoingWebhookMutation, article, showMessage, showErrorMessage]
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
