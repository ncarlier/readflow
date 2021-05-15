
const realm = process.env.KEYCLOAK_REALM || 'https://login.readflow.app/auth/admin/realms/readflow'
const clientID = process.env.KEYCLOAK_CLIENT_ID || 'subscription-management'
const clientSecret = process.env.KEYCLOAK_CLIENT_SECRET

const config = {
  realm,
  clientID,
  clientSecret
}

export default config
