/* eslint-disable @typescript-eslint/camelcase */
import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateOutboundServiceResponse, GetOutboundServicesResponse } from './models'
import { GetOutboundServices } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateOutboundServiceResponse | null }
) => {
  if (!mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateOutboundService
  const previousData = proxy.readQuery<GetOutboundServicesResponse>({
    query: GetOutboundServices,
  })
  if (previousData) {
    if (created.is_default) {
      previousData.outboundServices = previousData.outboundServices.map((service) => {
        return { ...service, is_default: false }
      })
    }
    const outboundServices = [created, ...previousData.outboundServices]
    proxy.writeQuery<GetOutboundServicesResponse>({ data: { outboundServices }, query: GetOutboundServices })
  }
}
