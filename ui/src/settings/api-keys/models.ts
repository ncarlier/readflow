export interface ApiKey {
  id: number
  alias: string
  token: string
  last_usage_at: string
  created_at?: string
  updated_at?: string
}

export interface GetApiKeysResponse {
  apiKeys: ApiKey[]
}

export interface GetApiKeyResponse {
  apiKey: ApiKey
}

export interface CreateOrUpdateApiKeyResponse {
  createOrUpdateAPIKey: ApiKey
}
