// API base URL
export const API_BASE_URL = __READFLOW_CONFIG__.apiBaseUrl || process.env.REACT_APP_API_ROOT || ''

// OIDC authority URL
export const AUTHORITY = __READFLOW_CONFIG__.authority || process.env.REACT_APP_AUTHORITY || 'none'

// OIDC client ID
export const CLIENT_ID = process.env.REACT_APP_CLIENT_ID || '232148523175444487@readflow.app'

// Unauthenticated user redirect
export const REDIRECT_URL = process.env.REACT_APP_REDIRECT_URL || '/login'
