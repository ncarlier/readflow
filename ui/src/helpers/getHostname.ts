export function getHostname(_url: string) {
  try {
    return new URL(_url).hostname
  } catch (e) {
    return _url
  }
}
