
import keycloak from '@/config/keycloak'
import oidcConfig from '@/config/oidc'

const getAccessToken = async () => {
  const res = await fetch(`${oidcConfig.authority}/protocol/openid-connect/token`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
    },
    body: new URLSearchParams({
      'grant_type': 'client_credentials',
      'client_id': keycloak.clientID,
      'client_secret': keycloak.clientSecret
    })
  })
  if (res.error) {
    throw res.error
  }
  const data = await res.json()
  if (!res.ok) {
    console.error('unable to get access token:', res.statusText)
    throw data.error
  }
  return data.access_token
}

export const updateUserAttributes = async (id, attributes) => {
  const token = await getAccessToken()
  console.debug('updating user with attributes', id, attributes)
  const res = await fetch(`${keycloak.realm}/users/${id}`, {
    method: 'PUT',
    headers: new Headers({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }),
    body: JSON.stringify({attributes})
  })

  if (res.error) {
    throw error
  }
  if (!res.ok) {
    console.error('unable to update user attributes:', res.statusText)
    throw res.statusText
  }
  console.info('user updated with attributes:', id, attributes)
  return
}

export const getUser = async (id) => {
  const token = await getAccessToken()
  const res = await fetch(`${keycloak.realm}/users/${id}`, {
    method: 'GET',
    headers: new Headers({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    })
  })

  if (res.error) {
    throw error
  }
  const data = await res.json()
  if (!res.ok) {
    console.error('unable to get user:', res.statusText)
    throw data.error
  }
  return data
}
