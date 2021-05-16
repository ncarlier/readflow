
const endpoint = process.env.READFLOW_ENDPOINT || 'https://api.readflow.app'
const clientID = process.env.READFLOW_CLIENT_ID || 'subscription-management'
const clientSecret = process.env.READFLOW_CLIENT_SECRET

const config = {
  endpoint,
  clientID,
  clientSecret
}

export default config
