import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateApiKeyResponse, GetApiKeysResponse } from './models'
import { GetApiKey, GetApiKeys } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateApiKeyResponse | null }
) => {
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys
  })
  if (previousData && mutationResult.data) {
    previousData.apiKeys.unshift(mutationResult.data.createOrUpdateAPIKey)
    proxy.writeQuery({ data: previousData, query: GetApiKeys })
  }
}

export const updateCacheAfterUpdate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateApiKeyResponse | null }
) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const updated = mutationResult.data.createOrUpdateAPIKey
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys
  })
  if (previousData) {
    const apiKeys = previousData.apiKeys.map(apiKey => {
      return apiKey.id === updated.id ? updated : apiKey
    })
    proxy.writeQuery({ data: { apiKeys }, query: GetApiKeys })
  }
  proxy.writeQuery({
    data: {
      apiKey: updated
    },
    query: GetApiKey,
    variables: { id: updated.id }
  })
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys
  })
  if (previousData) {
    const apiKeys = previousData.apiKeys.filter(apiKey => !ids.includes(apiKey.id))
    proxy.writeQuery({ data: { apiKeys }, query: GetApiKeys })
  }
}
