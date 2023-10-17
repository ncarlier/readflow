
const authParams = ['code', 'state', 'session_state', 'error'] as const

export const clearAuthParams = (params: URLSearchParams): string => {
  authParams.forEach(param => params.delete(param))
  return params.toString()
}

export const hasAuthParams = (params: URLSearchParams): boolean => {
  for (const param of authParams) {
    if (params.has(param)) {
      return true
    }
  }
  return false
}
