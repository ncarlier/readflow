// readflow UI runtime configuration
const __READFLOW_CONFIG__ = {
    // API base URL, default if empty
    // Values: URL (ex: `https://api.readflow.ap`)
    // Default: ''
    // Default can be overridden by setting ${REACT_APP_API_ROOT} env variable during build time
    apiBaseUrl: '',
    // Authorithy, default if empty
    // Values: URL if using OIDC (ex: `https://login.nunux.org/auth/realms/readflow`), `none` otherwise
    // Default: `none`
    // Default can be overridden by setting ${REACT_APP_AUTHORITY} env variable during build time
    authority: ''
}