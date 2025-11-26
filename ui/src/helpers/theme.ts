const getMql = () => window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)')

export const getEffectiveTheme = (theme: string) => {
  if (theme === 'auto') {
    const mql = getMql()
    return mql && mql.matches ? 'dark' : 'light'
  }
  return theme
}
