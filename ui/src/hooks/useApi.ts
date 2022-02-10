import { useEffect, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import { fetchAPI, withCredentials } from '../helpers'

const defaultHeaders = new Headers({
  'Content-Type': 'application/json',
})

export const useAPI = <T>(
  uri = '/',
  params: any = {},
  init: RequestInit = { headers: defaultHeaders }
): [boolean, T?, Error?] => {
  const { user } = useAuth()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error>()
  const [data, setData] = useState<T>()

  const stringifiedParams = JSON.stringify(params)
  useEffect(() => {
    const abortController = new AbortController()
    const headers = withCredentials(user, init.headers)
    const doFetchAPI = async () => {
      try {
        const res = await fetchAPI(uri, params, { ...init, signal: abortController.signal, headers })

        if (res.status >= 300) {
          throw new Error(res.statusText)
        }

        const data = await res.json()
        setData(data)
      } catch (e) {
        if (e.name !== 'AbortError') setError(e)
      } finally {
        setLoading(false)
      }
    }
    doFetchAPI()
    return () => abortController.abort()
    // eslint-disable-next-line
  }, [user, uri, stringifiedParams])

  return [loading, data, error]
}
