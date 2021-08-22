import React, { useCallback, useContext } from 'react'
import { useMutation, useQuery } from '@apollo/client'

import { Loader } from '../../../components'
import { getGQLError, matchResponse } from '../../../helpers'
import { GetOutgoingWebhooksResponse } from '../../../settings/intergrations/outgoing-webhook/models'
import { GetOutgoingWebhooks } from '../../../settings/intergrations/outgoing-webhook/queries'
import { Article } from '../../models'
import { SendArticleToOutgoingWebhook } from '../../queries'
import DefaultWebhookLink from './DefaultWebhookLink'
import OtherWebhooksLink from './OtherWebhooksLink'
import { MessageContext } from '../../../contexts/MessageContext'

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
  const { showMessage, showErrorMessage } = useContext(MessageContext)
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
      <>
        <li>
          <DefaultWebhookLink
            webhook={outgoingWebhooks.find((webhook) => webhook.is_default)}
            keyboard={keyboard}
            sendArticle={sendArticle}
          />
        </li>
        <li>
          <OtherWebhooksLink webhooks={outgoingWebhooks} sendArticle={sendArticle} />
        </li>
      </>
    ),
  })

  return <>{render(loading, data, error)}</>
}
