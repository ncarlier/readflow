
import api from '@/config/api'
import oidcConfig from '@/config/oidc'

const getRegisterUserMutation = (username) => `
mutation {
  registerUser(username: "${username}") {
    id
    username
    customer_id
    plan
  }
}
`
const getUpdateUserMutation = (id, payload) => {
  let params = ''
  Object.entries(payload).forEach(([key, value]) => {
    params += `,${key}: "${value}"`
  })
  return `
mutation {
  updateUser(uid: "${id}" ${params}) {
    id
    username
    customer_id
    plan
  }
}`
}

/**
 * Get OIDC access token.
 * @returns {string} OIDC access token
 */
const getAccessToken = async () => {
  const res = await fetch(`${oidcConfig.authority}/protocol/openid-connect/token`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
    },
    body: new URLSearchParams({
      'grant_type': 'client_credentials',
      'client_id': api.clientID,
      'client_secret': api.clientSecret
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

/**
 * Update readflow user.
 * @param {number} id user ID
 * @param {Object} payload update payload
 * @returns {Object} updated user
 */
export const updateUser = async (id, payload) => {
  const token = await getAccessToken()
  console.debug('updating user:', id, payload)
  const gql = {
    query: getUpdateUserMutation(id, payload),
    variables: null
  }
  const res = await fetch(`${api.endpoint}/admin`, {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }),
    body: JSON.stringify(gql)
  })

  if (res.error) {
    throw error
  }
  const data = await res.json()
  if (!res.ok) {
    console.error('unable to update user:', res.statusText)
    throw data.error
  }
  console.info('user updated:', data)
  return data.data.updateUser
}

/**
 * Get or register readflow user.
 * @param {string} username 
 */
export const getOrRegisterUser = async (username) => {
  const token = await getAccessToken()
  const gql = {
    query: getRegisterUserMutation(username),
    variables: null
  }
  const res = await fetch(`${api.endpoint}/admin`, {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }),
    body: JSON.stringify(gql)
  })

  if (res.error) {
    throw error
  }
  const data = await res.json()
  if (!res.ok) {
    console.error('unable to get or register user:', res.statusText)
    throw data.error
  }
  return data.data.registerUser
}
