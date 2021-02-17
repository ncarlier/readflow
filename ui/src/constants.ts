// API base URL
export const API_BASE_URL = process.env.REACT_APP_API_ROOT || 'https://api.readflow.app'

// OIDC authority URL
export const AUTHORITY = process.env.REACT_APP_AUTHORITY || 'https://login.readflow.app/auth/realms/readflow'

// OIDC client ID
export const CLIENT_ID = process.env.REACT_APP_CLIENT_ID || 'webapp'

// Unauthenticated user redirect
export const REDIRECT_URL = process.env.REACT_APP_REDIRECT_URL || 'https://about.readflow.app'

// VERSION
export const VERSION = process.env.REACT_APP_VERSION || 'snapshot'
