export function getOnlineStatus() {
  return typeof navigator !== 'undefined' && typeof navigator.onLine === 'boolean' ? navigator.onLine : true
}
