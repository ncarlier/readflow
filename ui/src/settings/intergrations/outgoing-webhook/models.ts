export interface OutgoingWebhook {
  id: number
  alias: string
  provider: string
  config: string
  is_default: boolean
  created_at?: string
  updated_at?: string
}

export interface GetOutgoingWebhooksResponse {
  outgoingWebhooks: OutgoingWebhook[]
}

export interface GetOutgoingWebhookResponse {
  outgoingWebhook: OutgoingWebhook
}

export interface CreateOrUpdateOutgoingWebhookRequest {
  id?: number
  alias: string
  provider: string
  config: string
  is_default: boolean
}

export interface CreateOrUpdateOutgoingWebhookResponse {
  createOrUpdateOutgoingWebhook: OutgoingWebhook
}

export interface DeleteOutgoingWebhookRequest {
  ids: number[]
}

export interface DeleteOutgoingWebhookResponse {
  deleteOutgoingWebhooks: number
}
