

export type Device = {
  id: number
  key: string
  created_at?: string
}

export interface GetDeviceResponse {
  device: Device
}

export type CreatePushSubscriptionResponse = {
  createPushSubscription: Device
}

export type DeletePushSubscriptionResponse = {
  deletePushSubscription: Device
}
