import React, { FC, PropsWithChildren, useCallback, useEffect, useState } from 'react'
import { useApolloClient, useMutation } from '@apollo/client'
import Switch from 'react-switch'

import { Box, Button, ErrorPanel, Loader } from '../../components'
import { isNotificationGranted, isNotificationSupported, subscribePush, unSubscribePush } from '../../helpers'
import { CreatePushSubscriptionResponse, DeletePushSubscriptionResponse, GetDeviceResponse } from './models'
import { CreatePushSubscription, DeletePushSubscription, GetDevice } from './queries'

const deviceIdKey = 'readflow.deviceId'

const NotificationSupport: FC<PropsWithChildren> = ({ children }) => {
  const supported = isNotificationSupported()
  const [allowed, setAllowed] = useState(isNotificationGranted())

  const requestPermission = () => Notification.requestPermission((permission) => setAllowed(permission === 'granted'))

  if (!supported) {
    return <p>Sorry, but this browser does not support desktop notification.</p>
  } else if (!allowed) {
    return (
      <>
        <p>Notifications are not yet allowed on your Browser.</p>
        <Button onClick={requestPermission}>Allow notifications</Button>
      </>
    )
  } else {
    return <>{children}</>
  }
}

interface NotificationErrorProps {
  reset: () => void
  err: Error
}

const NotificationError = ({ reset, err }: NotificationErrorProps) => (
  <ErrorPanel
    title="Unable to get push subcription status"
    actions={<Button onClick={reset}>reset subscription</Button>}
  >
    {err.message}
  </ErrorPanel>
)

const NotificationSwitch = () => {
  const [activated, setActivated] = useState(false)
  const [pushID, setPushID] = useState(localStorage.getItem(deviceIdKey))
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const client = useApolloClient()
  const [deletePushSubscriptionMutation] = useMutation<DeletePushSubscriptionResponse>(DeletePushSubscription)
  const [createPushSubscriptionMutation] = useMutation<CreatePushSubscriptionResponse>(CreatePushSubscription)

  const resetSubscription = async () => {
    localStorage.removeItem(deviceIdKey)
    setPushID(null)
    setActivated(false)
    try {
      const swr = await navigator.serviceWorker.ready
      await unSubscribePush(swr)
    } catch (err) {
      console.error(err)
    }
  }

  const getPushSubscriptionStatus = useCallback(
    async (pushId: string) => {
      setLoading(true)
      try {
        const { errors } = await client.query<GetDeviceResponse>({
          query: GetDevice,
          variables: { id: pushId },
        })
        if (errors) {
          throw new Error(errors[0].message)
        } else {
          setActivated(true)
        }
      } catch (err: any) {
        setError(err)
      } finally {
        setLoading(false)
      }
    },
    [client]
  )

  const subscribe = useCallback(async () => {
    try {
      setLoading(true)
      const swr = await navigator.serviceWorker.ready
      const subscription = await subscribePush(swr)
      if (subscription) {
        const res = await createPushSubscriptionMutation({
          variables: {
            sub: JSON.stringify(subscription),
          },
        })
        if (res.data) {
          const _id = res.data.createPushSubscription.id
          setPushID(_id.toString())
          localStorage.setItem(deviceIdKey, _id.toString())
        }
      }
    } catch (err: any) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [createPushSubscriptionMutation])

  const unsubscribe = useCallback(async () => {
    try {
      setLoading(true)
      await deletePushSubscriptionMutation({ variables: { id: pushID } })
      resetSubscription()
    } catch (err: any) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [deletePushSubscriptionMutation, pushID])

  useEffect(() => {
    if (pushID) {
      getPushSubscriptionStatus(pushID)
    }
  }, [pushID, getPushSubscriptionStatus])

  return (
    <>
      {error != null && <NotificationError err={error} reset={resetSubscription} />}
      <Switch onChange={activated ? unsubscribe : subscribe} checked={activated} />
      {loading && <Loader />}
    </>
  )
}

const NotificationBox = () => (
  <Box title="Notifications">
    <p>
      Receive notifications on this device.
    </p>
    <NotificationSupport>
      <NotificationSwitch />
    </NotificationSupport>
  </Box>
)

export default NotificationBox
