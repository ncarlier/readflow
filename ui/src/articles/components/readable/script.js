/* prettier-ignore */
window.onload = function() {
  document.querySelectorAll('img').forEach(i => ['src', 'srcset'].forEach(attr => i[attr] = i[attr].replaceAll('http://', 'https://')))
}
