import { UserManagerSettings } from 'oidc-client-ts'
import { AUTHORITY, CLIENT_ID } from '../config'

export const config: UserManagerSettings = {
  authority: AUTHORITY,
  client_id: CLIENT_ID,
  redirect_uri: `${document.location.origin}/login`,
  monitorSession: document.location.hostname !== 'localhost',
  response_type: 'code',
  scope: 'openid',
}
