import { useEffect, useContext } from 'react'
import { LocalConfigurationContext, Theme } from '../context/LocalConfigurationContext'

const getMql = () => window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)')

const getBrowserTheme = () => {
  const mql = getMql()
  return mql && mql.matches ? 'dark' : 'light'
}

const onBrowserThemeChanged = (callback: (theme: Theme) => void) => {
  const mql = getMql()
  const mqlListener = (e: MediaQueryListEvent) => callback(e.matches ? 'dark' : 'light')
  mql && mql.addListener(mqlListener)
  return () => mql && mql.removeListener(mqlListener)
}

export default () => {
  const { localConfiguration } = useContext(LocalConfigurationContext)

  const applyTheme = (theme: Theme) => {
    document.body.setAttribute('data-theme', theme)
  }

  useEffect(() => {
    const { theme } = localConfiguration
    console.log(`applying ${theme} theme`)
    if (theme === 'auto') {
      applyTheme(getBrowserTheme())
      return onBrowserThemeChanged(applyTheme)
    } else {
      applyTheme(theme)
    }
  }, [localConfiguration])
}
