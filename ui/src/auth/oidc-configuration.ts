import { UserManagerSettings } from 'oidc-client-ts'
import { AUTHORITY, CLIENT_ID } from '../constants'

const getRedirectURI = (): string => {
  const url = new URL(document.location.toString())
  url.searchParams.delete('code')
  url.searchParams.delete('state')
  url.searchParams.delete('session_state')
  url.searchParams.delete('error')
  console.debug('computed redirect URI:', url.toString())
  return url.toString()
}

export const config: UserManagerSettings = {
  authority: AUTHORITY,
  client_id: CLIENT_ID,
  redirect_uri: getRedirectURI(),
  monitorSession: document.location.hostname !== 'localhost',
  response_type: 'code',
  scope: 'openid',
}
