import { AUTHORITY, CLIENT_ID } from '../constants'

export const config = {
  authority: AUTHORITY,
  client_id: CLIENT_ID,
  redirect_uri: document.location.toString(),
  post_logout_redirect_uri: document.location.origin,
  monitorSession: true,
  revokeTokensOnSignout: true,
  response_type: 'code',
  scope: 'openid profile email',
}
