import { useApolloClient } from '@apollo/client'
import { useCallback, useEffect, useState } from 'react'
import { Device, GetDeviceResponse } from '../settings/preferences/models'
import { GetDevice } from '../settings/preferences/queries'

const deviceIdKey = 'readflow.deviceId'

export const useDeviceSubscriptionStatus = () => {
  const [device, setDevice] = useState<Device | null>(null)
  const deviceID = localStorage.getItem(deviceIdKey)
  const client = useApolloClient()

  const getPushSubscriptionStatus = useCallback(
    async (id: string) => {
      try {
        const {
          errors,
          data: { device },
        } = await client.query<GetDeviceResponse>({
          query: GetDevice,
          variables: { id },
        })
        if (errors) {
          throw new Error(errors[0].message)
        }
        setDevice(device)
      } catch (err: any) {
        console.error(err)
      }
    },
    [client]
  )

  useEffect(() => {
    if (deviceID) {
      getPushSubscriptionStatus(deviceID)
    }
  }, [getPushSubscriptionStatus])

  return device
}
