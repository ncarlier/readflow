// readflow UI runtime configuration
const __READFLOW_CONFIG__ = {
    // API base URL, default if empty
    // Values: URL (ex: `https://api.readflow.ap`)
    // Default: ''
    // Default can be overridden by setting ${REACT_APP_API_ROOT} env variable during build time
    apiBaseUrl: '{{ .HTTP.PublicURL }}',
    // Authorithy, default if empty
    // Values: URL if using OIDC (ex: `https://accounts.readflow.app`), `none` otherwise
    // Default: `none`
    // Default can be overridden by setting ${REACT_APP_AUTHORITY} env variable during build time
    authority: '{{ if eq .AuthN.Method "oidc" }}{{ .AuthN.OIDC.Issuer }}{{ else }}none{{ end }}',
    // OpenID Connect client ID, default if empty
    // Values: string (ex: `232148523175444487@readflow.app`)
    // Default: ''
    // Default can be overridden by setting ${REACT_APP_CLIENT_ID} env variable during build time
    client_id: '{{ .UI.ClientID }}',
}