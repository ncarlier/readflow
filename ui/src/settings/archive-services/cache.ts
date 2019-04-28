/* eslint-disable @typescript-eslint/camelcase */
import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateArchiveServiceResponse, GetArchiveServicesResponse } from './models'
import { GetArchiveService, GetArchiveServices } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data: CreateOrUpdateArchiveServiceResponse }
) => {
  const created = mutationResult.data.createOrUpdateArchiver
  const previousData = proxy.readQuery<GetArchiveServicesResponse>({
    query: GetArchiveServices
  })

  if (created.is_default && previousData) {
    previousData.archivers = previousData.archivers.map(service => {
      return { ...service, is_default: false }
    })
  }
  if (previousData) {
    previousData.archivers.unshift(created)
  }
  proxy.writeQuery({ data: previousData, query: GetArchiveServices })
}

export const updateCacheAfterUpdate = (
  proxy: DataProxy,
  mutationResult: { data: CreateOrUpdateArchiveServiceResponse }
) => {
  if (!mutationResult) {
    return
  }
  const updated = mutationResult.data.createOrUpdateArchiver
  const previousData = proxy.readQuery<GetArchiveServicesResponse>({
    query: GetArchiveServices
  })
  if (previousData) {
    const archivers = previousData.archivers.map(service => {
      if (updated.is_default) {
        service = { ...service, is_default: false }
      }
      return service.id === updated.id ? updated : service
    })
    proxy.writeQuery({ data: { archivers }, query: GetArchiveServices })
  }
  proxy.writeQuery({
    data: {
      archiver: updated
    },
    query: GetArchiveService,
    variables: { id: updated.id }
  })
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetArchiveServicesResponse>({
    query: GetArchiveServices
  })
  if (previousData) {
    const archivers = previousData.archivers.filter(archiver => archiver.id && !ids.includes(archiver.id))
    proxy.writeQuery({ data: { archivers }, query: GetArchiveServices })
  }
}
