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
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    return () => {
      if (loading && !abortController.abort) {
        abortController.abort()
      }
    }
  }, [abortController, loading])

  const doRequest = useCallback(
    async (params: any = {}) => {
      setLoading(true)
      const headers = withCredentials(user, init.headers)
      try {
        const res = await fetchAPI(uri, params, { ...init, signal: abortController.signal, headers })
        if (res.status >= 300) {
          throw new Error(res.statusText)
        }
        return res
      } catch (e: any) {
        console.error(e)
        if (e.name !== 'AbortError') throw e
      } finally {
        setLoading(false)
      }
    },
    [user, uri, init, abortController]
  )

  return doRequest
}
