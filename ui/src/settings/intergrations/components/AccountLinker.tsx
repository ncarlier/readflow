import React, { useCallback, useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'

import { Button, Loader } from '../../../components'
import { useMessage } from '../../../contexts'
import { fetchAPI } from '../../../helpers'

const sessionStoragePrefix = 'readflow.accountLinker'

const storeData = (provider: string, data: any) => {
  const key = sessionStoragePrefix + '.' + provider
  if (data === null) {
    window.sessionStorage.removeItem(key)
  } else {
    window.sessionStorage.setItem(key, JSON.stringify(data))
  }
}

const fetchData = (provider: string) => {
  const key = sessionStoragePrefix + '.' + provider
  const value = window.sessionStorage.getItem(key)
  if (!value) {
    return null
  }
  try {
    return JSON.parse(value)
  } catch {
    window.sessionStorage.removeItem(key)
  }
}

interface Props {
  provider: 'pocket'
}

export const AccountLinker = ({ provider }: Props) => {
  const [loading, setLoading] = useState(false)
  const { showErrorMessage } = useMessage()
  const history = useHistory()

  const linkAuthorize = useCallback(
    async (code: string) => {
      setLoading(true)
      let data: any
      try {
        const res = await fetchAPI(`/linking/${provider}/authorize`, { code }, { method: 'GET' })
        if (res.ok) {
          data = await res.json()
        } else {
          throw new Error(res.statusText)
        }
      } catch (err: any) {
        console.error(err.message)
      } finally {
        storeData(provider, null)
        setLoading(false)
      }
      if (data) {
        const qs = new URLSearchParams(data)
        qs.set('provider', provider)
        history.replace({
          search: qs.toString(),
        })
      }
    },
    [provider, history]
  )

  const linkRequest = useCallback(async () => {
    const qs = new URLSearchParams(window.location.search)
    qs.set('provider', provider)
    setLoading(true)
    let data: any
    try {
      const params = {
        redirect_uri: document.location.origin + document.location.pathname + '?' + qs.toString(),
      }
      const res = await fetchAPI(`/linking/${provider}/request`, params, { method: 'GET' })
      if (res.ok) {
        data = await res.json()
        storeData(provider, data)
      } else {
        const err = await res.json()
        throw new Error(err.detail || res.statusText)
      }
    } catch (err: any) {
      showErrorMessage(err.message)
      setLoading(false)
    }
    if (data && data['redirect']) {
      window.location.replace(data.redirect)
    }
  }, [provider, showErrorMessage])

  useEffect(() => {
    const data = fetchData(provider)
    if (data && data.code && !loading) {
      linkAuthorize(data.code)
    }
  }, [loading, provider, linkAuthorize])

  if (loading) {
    return <Loader center />
  }
  return (
    <Button title="Link your account..." onClick={linkRequest} icon="security">
      Link with {provider}
    </Button>
  )
}
