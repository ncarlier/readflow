/* eslint-disable @typescript-eslint/camelcase */
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
    if (created.is_default) {
      previousData.outgoingWebhooks = previousData.outgoingWebhooks.map((service) => {
        return { ...service, is_default: false }
      })
    }
    const outgoingWebhooks = [created, ...previousData.outgoingWebhooks]
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
