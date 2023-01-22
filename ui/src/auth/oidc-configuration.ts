import { UserManagerSettings } from 'oidc-client-ts'
import { AUTHORITY, CLIENT_ID } from '../constants'
import { getCleanedRedirectURI } from './helper'

export const config: UserManagerSettings = {
  authority: AUTHORITY,
  client_id: CLIENT_ID,
  redirect_uri: getCleanedRedirectURI(document.location.href),
  monitorSession: document.location.hostname !== 'localhost',
  response_type: 'code',
  scope: 'openid',
}
