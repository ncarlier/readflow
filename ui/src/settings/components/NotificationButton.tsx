import React, { useEffect, useState } from 'react'
import { useApolloClient, useMutation } from 'react-apollo-hooks'

import ButtonIcon from '../../common/ButtonIcon'
import { CreatePushSubscriptionResponse, DeletePushSubscriptionResponse, GetDeviceResponse } from './models'
import { CreatePushSubscription, DeletePushSubscription, GetDevice } from './queries'

const DEVICE_ID = 'device_id'

export default () => {
  const [id, setId] = useState(localStorage.getItem(DEVICE_ID))
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const client = useApolloClient()

  const deletePushSubscriptionMutation = useMutation<DeletePushSubscriptionResponse>(DeletePushSubscription)
  const deletePushSubscription = async () => {
    try {
      setLoading(true)
      await deletePushSubscriptionMutation({ variables: { id } })
      setId(null)
      localStorage.removeItem(DEVICE_ID)
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }

  const getPushSubscription = async (pushId: string) => {
    setLoading(true)
    try {
      const { errors } = await client.query<GetDeviceResponse>({
        query: GetDevice,
        variables: { id: pushId }
      })
      if (errors) {
        throw new Error(errors[0])
      }
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (id) {
      getPushSubscription(id)
    }
  }, [id])

  const createPushSubscriptionMutation = useMutation<CreatePushSubscriptionResponse>(CreatePushSubscription)
  const createPushSubscription = async () => {
    try {
      setLoading(true)
      const swr = await navigator.serviceWorker.ready
      const subscription = await swr.pushManager.getSubscription()
      if (subscription) {
        const res = await createPushSubscriptionMutation({
          variables: {
            sub: JSON.stringify(subscription)
          }
        })
        const _id = res.data.createPushSubscription.id
        setId(_id)
        localStorage.setItem(DEVICE_ID, _id)
      }
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }

  const resetSubscription = async (err: Error) => {
    if (confirm(`An error occured:\n${err.message}\n\nReset subscription?`)) {
      await deletePushSubscription()
    }
  }

  let title = id ? 'Disable notifications' : 'Enable notifications'
  let icon = id ? 'notifications_off' : 'notifications'
  let onClick = id ? deletePushSubscription : createPushSubscription
  if (error) {
    title = error.message
    icon = 'notification_important'
    onClick = () => resetSubscription(error)
  }

  return <ButtonIcon title={title} onClick={onClick} loading={loading} icon={icon} />
}
