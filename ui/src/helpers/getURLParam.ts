export function getURLParam<T>(params: URLSearchParams, name: string, fallback: T): T {
  if (params.has(name)) {
    const val = params.get(name)
    if (val && typeof fallback === 'number') {
      // eslint-disable-next-line use-isnan
      if (!Number.isNaN(parseInt(val, 10))) {
        return Number.parseInt(val, 10) as any
      }
    }
    if (val && typeof fallback === 'string') {
      return val as any
    }
  }
  return fallback
}
