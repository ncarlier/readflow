import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateIncomingWebhookResponse, GetIncomingWebhooksResponse } from './models'
import { GetIncomingWebhooks } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateIncomingWebhookResponse | null }
) => {
  if (!mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateIncomingWebhook
  const previousData = proxy.readQuery<GetIncomingWebhooksResponse>({
    query: GetIncomingWebhooks,
  })
  if (previousData) {
    const incomingWebhooks = [created, ...previousData.incomingWebhooks]
    proxy.writeQuery<GetIncomingWebhooksResponse>({ data: { incomingWebhooks }, query: GetIncomingWebhooks })
  }
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetIncomingWebhooksResponse>({
    query: GetIncomingWebhooks,
  })
  if (previousData) {
    const incomingWebhooks = previousData.incomingWebhooks.filter((webhook) => !ids.includes(webhook.id))
    proxy.writeQuery<GetIncomingWebhooksResponse>({ data: { incomingWebhooks }, query: GetIncomingWebhooks })
  }
}
