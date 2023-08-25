export interface IncomingWebhook {
  id: number
  alias: string
  token: string
  email: string
  script: string
  last_usage_at: string
  created_at?: string
  updated_at?: string
}

export interface GetIncomingWebhooksResponse {
  incomingWebhooks: IncomingWebhook[]
}

export interface GetIncomingWebhookResponse {
  incomingWebhook: IncomingWebhook
}

export interface CreateOrUpdateIncomingWebhookRequest {
  id?: number
  alias: string
  script: string
}

export interface CreateOrUpdateIncomingWebhookResponse {
  createOrUpdateIncomingWebhook: IncomingWebhook
}

export interface DeleteIncomingWebhookRequest {
  ids: number[]
}

export interface DeleteIncomingWebhookResponse {
  deleteIncomingWebhooks: number
}
