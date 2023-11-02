import { User } from 'oidc-client-ts'
import { API_BASE_URL } from '../config'

const apiBaseUrl = API_BASE_URL.startsWith('/') || API_BASE_URL === '' ? document.location.origin + API_BASE_URL : API_BASE_URL

export const getAPIURL = (path?: string) => path ? `${apiBaseUrl}${path}` : apiBaseUrl

export const withCredentials = (user: User | null, init?: HeadersInit): HeadersInit | undefined => {
  if (user && user.access_token) {
    const headers = new Headers({ Authorization: 'Bearer ' + user.access_token })
    if (init !== undefined) {
      ;(init as Headers).forEach((entry) => {
        headers.set(entry[0], entry[1])
      })
    }
    init = headers
  }
  return init
}

export const fetchAPI = async (uri: string, params: any = {}, init: RequestInit = {}) => {
  const url = new URL(getAPIURL(uri))
  if (params) {
    Object.keys(params).forEach((key) => url.searchParams.append(key, params[key]))
  }
  return await fetch(url.toString(), init)
}
