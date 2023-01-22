
const authParams = ['code', 'state', 'session_state', 'error'] as const

export const clearAuthParams = (params: URLSearchParams): URLSearchParams => {
  authParams.forEach(param => params.delete(param))
  return params
}

export const hasAuthParams = (params: URLSearchParams): boolean => {
  for (const param of authParams) {
    if (params.has(param)) {
      return true
    }
  }
  return false
}

export const getCleanedRedirectURI = (href: string): string => {
  const url = new URL(href)
  clearAuthParams(url.searchParams)
  console.debug('computed redirect URI:', url.href)
  return url.href
}
