import { DataProxy } from "apollo-cache"
import { GetApiKeysResponse, CreateOrUpdateApiKeyResponse } from "./models"
import { GetApiKeys, GetApiKey } from "./queries"

export const updateCacheAfterCreate = (proxy: DataProxy, mutationResult: {data: CreateOrUpdateApiKeyResponse}) => {
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys,
  })
  previousData!.apiKeys.unshift(mutationResult!.data!.createOrUpdateAPIKey)
  proxy.writeQuery({ data: previousData, query: GetApiKeys })
}

export const updateCacheAfterUpdate = (proxy: DataProxy, mutationResult: {data: CreateOrUpdateApiKeyResponse}) => {
  const updated = mutationResult!.data.createOrUpdateAPIKey
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys,
  })
  const apiKeys = previousData!.apiKeys.map(apiKey => {
    return apiKey.id === updated.id ? updated : apiKey
  })
  proxy.writeQuery({ data: {apiKeys}, query: GetApiKeys })
  proxy.writeQuery({
    data: {
      apiKey: updated
    }, 
    query: GetApiKey,
    variables: {id: updated.id}
  })
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetApiKeysResponse>({
    query: GetApiKeys,
  })
  const apiKeys = previousData!.apiKeys.filter(apiKey => !ids.includes(apiKey.id))
  proxy.writeQuery({ data: {apiKeys}, query: GetApiKeys })
}
