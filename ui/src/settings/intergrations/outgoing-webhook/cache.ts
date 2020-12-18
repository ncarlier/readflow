/* eslint-disable @typescript-eslint/camelcase */
import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateOutgoingWebhookResponse, GetOutgoingWebhooksResponse } from './models'
import { GetOutgoingWebhooks } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateOutgoingWebhookResponse | null }
) => {
  if (!mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateOutgoingWebhook
  const previousData = proxy.readQuery<GetOutgoingWebhooksResponse>({
    query: GetOutgoingWebhooks,
  })
  if (previousData) {
    if (created.is_default) {
      previousData.outgoingWebhooks = previousData.outgoingWebhooks.map((service) => {
        return { ...service, is_default: false }
      })
    }
    const outgoingWebhooks = [created, ...previousData.outgoingWebhooks]
    proxy.writeQuery<GetOutgoingWebhooksResponse>({ data: { outgoingWebhooks }, query: GetOutgoingWebhooks })
  }
}
