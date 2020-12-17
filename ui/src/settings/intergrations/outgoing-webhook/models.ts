export interface OutboundService {
  id?: number
  alias: string
  provider: string
  config: string
  is_default: boolean
  created_at?: string
  updated_at?: string
}

export interface GetOutboundServicesResponse {
  outboundServices: OutboundService[]
}

export interface GetOutboundServiceResponse {
  outboundService: OutboundService
}

export interface CreateOrUpdateOutboundServiceResponse {
  createOrUpdateOutboundService: OutboundService
}

export interface DeleteOutboundServiceRequest {
  ids: number[]
}

export interface DeleteOutboundServiceResponse {
  deleteOutboundServices: number
}
