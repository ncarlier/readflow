import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateApiKeyResponse, GetApiKeysResponse } from './models'
import { GetApiKeys } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateApiKeyResponse | null }
) => {
  if (!mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateAPIKey
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys,
  })
  if (previousData) {
    const apiKeys = [created, ...previousData.apiKeys]
    proxy.writeQuery<GetApiKeysResponse>({ data: { apiKeys }, query: GetApiKeys })
  }
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys,
  })
  if (previousData) {
    const apiKeys = previousData.apiKeys.filter((apiKey) => !ids.includes(apiKey.id))
    proxy.writeQuery<GetApiKeysResponse>({ data: { apiKeys }, query: GetApiKeys })
  }
}
