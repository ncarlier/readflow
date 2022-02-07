import { API_BASE_URL } from '../constants'

export const withCredentials = (user?: any, init?: HeadersInit): HeadersInit | undefined => {
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

export const fetchAPI = async (uri: string, params: any = {}, init: RequestInit) => {
  const url = new URL(`${API_BASE_URL}${uri}`)
  if (params) {
    Object.keys(params).forEach((key) => url.searchParams.append(key, params[key]))
  }
  return await fetch(url.toString(), init)
}
