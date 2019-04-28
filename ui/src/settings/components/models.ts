export interface Device {
  id: number
  key: string
  created_at?: string
}

export interface GetDeviceResponse {
  device: Device
}

export interface CreatePushSubscriptionResponse {
  createPushSubscription: Device
}

export interface DeletePushSubscriptionResponse {
  deletePushSubscription: Device
}
