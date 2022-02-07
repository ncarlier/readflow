import { AUTHORITY, CLIENT_ID } from '../constants'

export const buildConfig = (redirect?: string) => {
  const params = new URLSearchParams()
  if (redirect) {
    params.append('redirect', redirect)
  }
  return {
    authority: AUTHORITY,
    client_id: CLIENT_ID,
    redirect_uri: `${document.location.origin}?${params.toString()}`,
    post_logout_redirect_uri: document.location.origin,
    monitorSession: true,
    revokeTokensOnSignout: true,
    response_type: 'code',
    scope: 'openid profile email',
  }
}
