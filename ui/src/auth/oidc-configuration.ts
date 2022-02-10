import { UserManagerSettings } from 'oidc-client-ts'
import { AUTHORITY, CLIENT_ID } from '../constants'

const getRedirectURL = (): string => {
  const url = new URL(document.location.toString())
  url.searchParams.delete('code')
  url.searchParams.delete('state')
  url.searchParams.delete('session_state')
  console.log('redirection', url.toString())
  return url.toString()
}

export const config: UserManagerSettings = {
  authority: AUTHORITY,
  client_id: CLIENT_ID,
  redirect_uri: getRedirectURL(),
  // post_logout_redirect_uri: document.location.origin + '/logout',
  monitorSession: true,
  response_type: 'code',
  scope: 'openid profile email',
}
