import { DataProxy } from '@apollo/client'

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
    let outgoingWebhooks = [created, ...previousData.outgoingWebhooks]
    if (created.is_default) {
      outgoingWebhooks = outgoingWebhooks.map((service) => service.id !== created.id ? { ...service, is_default: false } : service)
    }
    proxy.writeQuery<GetOutgoingWebhooksResponse>({ data: { outgoingWebhooks }, query: GetOutgoingWebhooks })
  }
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetOutgoingWebhooksResponse>({
    query: GetOutgoingWebhooks,
  })
  if (previousData) {
    const outgoingWebhooks = previousData.outgoingWebhooks.filter((webhook) => !ids.includes(webhook.id))
    proxy.writeQuery<GetOutgoingWebhooksResponse>({ data: { outgoingWebhooks }, query: GetOutgoingWebhooks })
  }
}
