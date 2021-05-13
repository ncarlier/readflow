
const authority = process.env.REACT_APP_AUTHORITY || 'https://login.readflow.app/auth/realms/readflow'
const clientId = process.env.REACT_APP_CLIENT_ID || 'about.readflow.app'
const redirectUri = process.env.NODE_ENV === "development" ? "http://localhost:3000/" : "https://about.readflow.app"

const config = {
  onSignIn: async (user) => {
    console.log(`logged user: ${user}`)
  },
  authority,
  clientId,
  responseType: "id_token",
  scope: "openid profile",
  autoSignIn: false,
  redirectUri
}

export default config
