// API base URL
export const API_BASE_URL = __READFLOW_CONFIG__.apiBaseUrl || process.env.REACT_APP_API_ROOT || ''

// OIDC authority URL
export const AUTHORITY = __READFLOW_CONFIG__.authority || process.env.REACT_APP_AUTHORITY || 'none'

// OIDC client ID
export const CLIENT_ID = __READFLOW_CONFIG__.client_id || process.env.REACT_APP_CLIENT_ID || ''

// Portal URL
export const PORTAL_URL = process.env.REACT_APP_PORTAL_URL || ''
