/* eslint-disable @typescript-eslint/camelcase */
import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateArchiveServiceResponse, GetArchiveServicesResponse } from './models'
import { GetArchiveServices } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateArchiveServiceResponse | null }
) => {
  if (!mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateArchiver
  const previousData = proxy.readQuery<GetArchiveServicesResponse>({
    query: GetArchiveServices,
  })
  if (previousData) {
    if (created.is_default) {
      previousData.archivers = previousData.archivers.map((service) => {
        return { ...service, is_default: false }
      })
    }
    const archivers = [created, ...previousData.archivers]
    proxy.writeQuery<GetArchiveServicesResponse>({ data: { archivers }, query: GetArchiveServices })
  }
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetArchiveServicesResponse>({
    query: GetArchiveServices,
  })
  if (previousData) {
    const archivers = previousData.archivers.filter((archiver) => archiver.id && !ids.includes(archiver.id))
    proxy.writeQuery({ data: { archivers }, query: GetArchiveServices })
  }
}
