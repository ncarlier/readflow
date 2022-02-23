import { useCallback, useEffect, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import { fetchAPI, withCredentials } from '../helpers'

const defaultHeaders = new Headers({
  'Content-Type': 'application/json',
})

type doRequestFn = (params?: any) => Promise<Response | undefined>

export const useAPI = (uri = '/', init: RequestInit = { headers: defaultHeaders }): doRequestFn => {
  const { user } = useAuth()
  const [abortController] = useState(() => new AbortController())

  useEffect(() => {
    return () => abortController.abort()
  }, [abortController])

  const doRequest = useCallback(
    async (params: any = {}) => {
      const headers = withCredentials(user, init.headers)
      try {
        const res = await fetchAPI(uri, params, { ...init, signal: abortController.signal, headers })
        if (res.status >= 300) {
          throw new Error(res.statusText)
        }
        return res
      } catch (e) {
        if (e.name !== 'AbortError') throw e
      }
    },
    [user, uri, init, abortController]
  )

  return doRequest
}
