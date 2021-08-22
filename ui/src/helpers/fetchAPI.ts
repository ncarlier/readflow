import authService from '../auth'
import { API_BASE_URL } from '../constants'

export const fetchAPI = async (uri: string, params: any = {}, init: RequestInit) => {
  let user = await authService.getUser()
  if (user === null) {
    throw new Error('user not logged in')
  }
  if (user.expired) {
    user = await authService.renewToken()
  }
  if (user.access_token) {
    const headers = new Headers({ Authorization: 'Bearer ' + user.access_token })
    if (init.headers !== undefined) {
      ;(init.headers as Headers).forEach((entry) => {
        headers.set(entry[0], entry[1])
      })
    }
    init.headers = headers
  }
  const url = new URL(`${API_BASE_URL}${uri}`)
  if (params) {
    Object.keys(params).forEach((key) => url.searchParams.append(key, params[key]))
  }
  return await fetch(url.toString(), init)
}
